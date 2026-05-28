package main

import (
	"encoding/binary"
)

type DNSQuestion struct {
	Name   string
	Type   uint16
	Class  uint16
	Offset int
}

func (q *DNSQuestion) MarshalBinary() (data []byte, err error) {
	data = []byte{}
	encoded, err := encodeDomainName(q.Name)
	if err != nil {
		return []byte{}, err
	}
	data = append(data, encoded...)
	data = append(data, 0x00)
	data = binary.BigEndian.AppendUint16(data, q.Type)
	data = binary.BigEndian.AppendUint16(data, q.Class)
	return data, nil
}

func (q *DNSQuestion) UnmarshalBinary(buf []byte) error {
	name, offset, err := decodeDomainName(buf, q.Offset)
	if err != nil {
		return err
	}

	q.Name = name
	q.Type = binary.BigEndian.Uint16(buf[offset : offset+2])
	q.Class = binary.BigEndian.Uint16(buf[offset+2 : offset+4])
	q.Offset = offset + 4

	return nil
}
