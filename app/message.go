package main

import (
	"encoding/binary"
)

type Message struct {
	Header Header
}

type Header struct {
	ID      uint16
	QR      uint8
	OPCODE  uint8
	AA      uint8
	TC      uint8
	RD      uint8
	RA      uint8
	Z       uint8
	RCODE   uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func NewMessage() *Message {
	return &Message{}
}

func NewHeader() *Header {
	return &Header{
		ID:      1234,
		QR:      1,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      0,
		RA:      0,
		Z:       0,
		RCODE:   0,
		QDCOUNT: 0,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
}

func (m *Message) Serialize() []byte {
	header := NewHeader()
	buf := make([]byte, 0, 12)
	flags := uint16(0)
	flags |= uint16(header.QR) << 15
	flags |= uint16(header.OPCODE) << 11
	flags |= uint16(header.AA) << 10
	flags |= uint16(header.TC) << 9
	flags |= uint16(header.RD) << 8
	flags |= uint16(header.RA) << 7
	flags |= uint16(header.Z) << 4
	flags |= uint16(header.RCODE)
	buf = binary.BigEndian.AppendUint16(buf, header.ID)
	buf = binary.BigEndian.AppendUint16(buf, flags)
	buf = binary.BigEndian.AppendUint16(buf, header.QDCOUNT)
	buf = binary.BigEndian.AppendUint16(buf, header.ANCOUNT)
	buf = binary.BigEndian.AppendUint16(buf, header.NSCOUNT)
	buf = binary.BigEndian.AppendUint16(buf, header.ARCOUNT)

	return buf
}
