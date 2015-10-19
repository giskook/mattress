package mattress

import (
	"github.com/giskook/gotcp"
	"log"
)

var (
	Illegal  uint16 = 0
	HalfPack uint16 = 255

	ReportStatus uint16 = 1
)

type MattressProtocol struct {
}

func (this *MattressProtocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	smconn := c.GetExtraData().(*Conn)
	buffer := smconn.GetBuffer()

	conn := c.GetRawConn()
	for {
		data := make([]byte, 2048)
		readLengh, err := conn.Read(data)
		log.Printf("%X\n", data[0:readLengh])

		if err != nil {
			return nil, err
		}

		if readLengh == 0 {
			return nil, gotcp.ErrConnClosing
		} else {
			buffer.Write(data[0:readLengh])
			cmdid, pkglen := CheckProtocol(buffer)
			//		log.Printf("recv box cmd %d \n", cmdid)

			pkgbyte := make([]byte, pkglen)
			buffer.Read(pkgbyte)
			switch cmdid {
			case ReportStatus:
				return ParseReportStatus(pkgbyte), nil
			case HalfPack:
			case Illegal:
			}
		}
	}
}
