package mattress

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

var (
	MaxLength      uint16 = 16
	Wet            uint8  = 0
	Dried          uint8  = 1
	NurseCall      uint8  = 0x30
	NurseCallTwice uint8  = 0x31
	BodyMovement   uint8  = 0x41
	LeaveBed       uint8  = 0x50
	InBed          uint8  = 0x20
)

func CheckSum(buffer []byte) uint8 {
	var sum uint16 = 0
	for i := 0; i < len(buffer); i++ {
		sum += uint16(buffer[i])
	}
	return uint8(sum)&0x7F + 0x80
}

func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	if buffer.Bytes()[0] != 0x06 {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else if uint16(bufferlen) < MaxLength {
		return HalfPack, 0
	} else if buffer.Bytes()[15] != 0x04 {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else {
		checksum := CheckSum(buffer.Bytes()[1:14])
		if checksum == buffer.Bytes()[14] {
			return ReportStatus, uint16(bufferlen)
		} else {
			CheckProtocol(buffer)
		}
	}

	return Illegal, 0
}

type ReportStatusPacket struct {
	MattressID uint64
	Status     uint8
	Heartbeat  uint16
	Breath     uint8
	TrunOver   uint8
}

func (p *ReportStatusPacket) Serialize() []byte {
	return nil
}

func ParseReportStatus(buffer []byte) *ReportStatusPacket {
	reader := bytes.NewReader(buffer)
	reader.ReadByte()
	mattressid_byte := make([]byte, 6)
	reader.Read(mattressid_byte)
	midbytes := []byte{0, 0}
	midbytes = append(midbytes, mattressid_byte...)
	mattressid := binary.BigEndian.Uint64(midbytes)
	status, _ := reader.ReadByte()
	heart_byte := make([]byte, 3)
	reader.Read(heart_byte)
	heart, _ := strconv.Atoi(string(heart_byte))
	breath_byte := make([]byte, 2)
	reader.Read(breath_byte)
	breath, _ := strconv.Atoi(string(breath_byte))
	turnover, _ := reader.ReadByte()

	return &ReportStatusPacket{
		MattressID: mattressid,
		Status:     status,
		Heartbeat:  uint16(heart),
		Breath:     uint8(breath),
		TrunOver:   turnover,
	}
}
