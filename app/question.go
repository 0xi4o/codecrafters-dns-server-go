package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type DNSQuestion struct {
	Name   string
	Type   uint16
	Class  uint16
	Offset int
}

func (q *DNSQuestion) MarshalBinary() (data []byte, err error) {
	data = []byte{}
	labels := strings.Split(q.Name, ".")
	for _, label := range labels {
		length := byte(len(label))
		data = append(data, length)
		data = append(data, []byte(label)...)
	}
	data = binary.BigEndian.AppendUint16(data, 0x00)
	data = binary.BigEndian.AppendUint16(data, q.Type)
	data = binary.BigEndian.AppendUint16(data, q.Class)
	return data, nil
}

func (q *DNSQuestion) UnmarshalBinary(buf []byte) error {
	name, offset, err := parseName(buf, 0)
	if err != nil {
		return err
	}

	q.Name = name
	q.Type = binary.BigEndian.Uint16(buf[offset : offset+2])
	q.Class = binary.BigEndian.Uint16(buf[offset+2 : offset+4])

	return nil
}

func parseName(buf []byte, offset int) (string, int, error) {
	labels := []string{}
	for {
		if offset >= len(buf) {
			return "", 0, fmt.Errorf("buffer too short for parsing name")
		}

		length := int(buf[offset])
		offset++

		if length == 0 {
			break
		}

		if offset+length > len(buf) {
			return "", 0, fmt.Errorf("label length exceeds buffer")
		}

		labels = append(labels, string(buf[offset:offset+length]))
		offset += length
	}
	return strings.Join(labels, "."), offset, nil
}
