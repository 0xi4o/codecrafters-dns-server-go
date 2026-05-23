package main

import (
	"fmt"
	"strconv"
	"strings"
)

func encodeARecord(ip string) ([]byte, error) {
	data := make([]byte, 4)
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return []byte{}, fmt.Errorf("invalid A record")
	}
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return []byte{}, fmt.Errorf("invalid part at index: %d", i)
		}
		data = append(data, uint8(num)&0x0F)
	}
	return data, nil
}

func encodeDomainName(name string) ([]byte, error) {
	data := []byte{}
	labels := strings.Split(name, ".")
	if len(labels) < 2 {
		return []byte{}, fmt.Errorf("invalid domain name")
	}
	for _, label := range labels {
		length := byte(len(label))
		data = append(data, length)
		data = append(data, []byte(label)...)
	}
	return data, nil
}

func decodeDomainName(buf []byte, offset int) (string, int, error) {
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
