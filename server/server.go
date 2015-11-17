package main

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/mattress"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// creates a mqtt client
	mqttclient := mattress.NewMqttClient()
	mqttclient.Connection()
	mattress.SetMqttClient(mqttclient)
	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":5858")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a tcp server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &mattress.Callback{}, &mattress.MattressProtocol{})

	go srv.Start(listener, time.Second)

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
