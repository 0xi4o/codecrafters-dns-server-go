package main

import (
	"encoding/binary"
	"fmt"
)

type DNSHeader struct {
	ID      uint16
	QR      bool
	OPCODE  uint8
	AA      bool
	TC      bool
	RD      bool
	RA      bool
	Z       uint8
	RCODE   uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h *DNSHeader) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 12)
	binary.BigEndian.PutUint16(data[:2], h.ID)
	var part1 byte
	var part2 byte
	// QR, OPCODE, AA, TC, and RD
	if h.QR {
		part1 |= 0b1000_0000
	}
	part1 |= (h.OPCODE & 0x0F) << 3
	if h.AA {
		part1 |= 0b0000_0100
	}
	if h.TC {
		part1 |= 0b0000_0010
	}
	if h.RD {
		part1 |= 0b0000_0001
	}
	// RA, Z, and RCODE
	if h.RA {
		part2 |= 0b1000_0000
	}
	part2 |= (h.Z & 0x07) << 4
	part2 |= h.RCODE & 0x0F

	data[2] = part1
	data[3] = part2
	binary.BigEndian.PutUint16(data[4:6], h.QDCOUNT)
	binary.BigEndian.PutUint16(data[6:8], h.ANCOUNT)
	binary.BigEndian.PutUint16(data[8:10], h.NSCOUNT)
	binary.BigEndian.PutUint16(data[10:12], h.ARCOUNT)
	return data, err
}

func (h *DNSHeader) UnmarshalBinary(buf []byte) error {
	if len(buf) != 12 {
		return fmt.Errorf("Length of DNS Header should be 12 bytes")
	}
	// part1 - QR (1 bit), OPCODE (4 bits), AA (1 bit), TC (1 bit), RD (1 bit)
	// part2 - RA (1 bit), Z (3 bits), and RCODE (4 bits)
	part1, part2 := buf[2], buf[3]
	h.ID = binary.BigEndian.Uint16(buf[:2])
	h.QR = part1&0b1000_0000 != 0
	h.OPCODE = (part1 >> 3) & 0x0F
	h.AA = part1&0b0000_0100 != 0
	h.TC = part1&0b0000_0010 != 0
	h.RD = part1&0b0000_0001 != 0
	h.RA = part2&0b1000_0000 != 0
	h.Z = (part2 >> 4) & 0x07
	h.RCODE = part2 & 0x0F
	h.QDCOUNT = binary.BigEndian.Uint16(buf[4:6])
	h.ANCOUNT = binary.BigEndian.Uint16(buf[6:8])
	h.NSCOUNT = binary.BigEndian.Uint16(buf[8:10])
	h.ARCOUNT = binary.BigEndian.Uint16(buf[10:12])
	return nil
}
