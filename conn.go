package mattress

import (
	"bytes"
	"encoding/json"
	"github.com/giskook/gotcp"
	"log"
)

type Conn struct {
	conn          *gotcp.Conn
	recieveBuffer *bytes.Buffer
	mattressid    uint64
}

func NewConn(conn *gotcp.Conn) *Conn {
	return &Conn{
		conn:          conn,
		recieveBuffer: bytes.NewBuffer([]byte{}),
	}
}

func (c *Conn) Close() {
	c.recieveBuffer.Reset()
}

func (c *Conn) GetBuffer() *bytes.Buffer {
	return c.recieveBuffer
}

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	conn := NewConn(c)

	c.PutExtraData(conn)

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	conn := c.GetExtraData().(*Conn)
	conn.Close()
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*ReportStatusPacket)
	res, _ := json.Marshal(packet)

	GetMqttClient().Publish("hello", string(res))
	log.Println(string(res))
	return true
}
