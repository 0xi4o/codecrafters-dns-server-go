package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

func DeserializeDomainOrIP(buf []byte, startingPosition int) (string, int) {
	var labels []string
	position := startingPosition
	seenPositions := make(map[int]bool)

	for position < len(buf) {
		if (buf[position] & 0xC0) == 0xC0 {
			if position+1 >= len(buf) {
				break
			}

			pointerOffset := int((uint16(buf[position]&0x3F)<<8 | uint16(buf[position+1])))

			returnPosition := position + 2

			if seenPositions[pointerOffset] {
				fmt.Println("Warning: Circular reference detected in DNS packet")
				break
			}
			seenPositions[pointerOffset] = true

			suffix, _ := DeserializeDomainOrIP(buf, pointerOffset)

			if len(labels) > 0 {
				return strings.Join(labels, ".") + "." + suffix, returnPosition
			}
			return suffix, returnPosition
		}

		labelLength := int(buf[position])
		position++

		if labelLength == 0 {
			break
		}

		if position+labelLength > len(buf) {
			fmt.Println("Warning: Label length exceeds buffer size")
			break
		}

		label := string(buf[position : position+labelLength])
		labels = append(labels, label)
		position += labelLength
	}

	if len(labels) > 0 {
		return strings.Join(labels, "."), position
	}
	return "", position
}

func SerializeDomainOrIP(domain string) []byte {
	if domain == "" {
		return []byte{0}
	}

	result := []byte{}
	labels := strings.Split(domain, ".")

	for _, label := range labels {
		if len(label) > 0 {
			result = append(result, byte(len(label)))
			result = append(result, []byte(label)...)
		}
	}

	result = append(result, 0)
	return result
}

func SerializeIPv4(ipStr string) []byte {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		fmt.Printf("Invalid IPv4 address: %s\n", ipStr)
		return []byte{0, 0, 0, 0}
	}
	return []byte(ip)
}

func ExtractAnswerFromResponse(response []byte, offset int) (Answer, error) {
	if offset >= len(response) {
		return Answer{}, fmt.Errorf("offset beyond response length")
	}

	name, nameOffset := DeserializeDomainOrIP(response, offset)
	offset = nameOffset

	if offset+10 > len(response) {
		return Answer{}, fmt.Errorf("response too short for answer fields")
	}

	ttl := binary.BigEndian.Uint32(response[offset+4 : offset+8])

	if offset+10+4 > len(response) {
		return Answer{}, fmt.Errorf("response too short for IPv4 address")
	}

	ipData := response[offset+10 : offset+14]
	ipStr := fmt.Sprintf("%d.%d.%d.%d", ipData[0], ipData[1], ipData[2], ipData[3])

	return Answer{
		Name:   name,
		Type:   1, // A record
		Class:  1, // IN class
		TTL:    ttl,
		Length: 4, // IPv4 address is always 4 bytes
		Data:   ipStr,
	}, nil
}

func CalculateAnswerSize(answer Answer) int {
	// Size for A records: name + type(2) + class(2) + TTL(4) + length(2) + IPv4(4)
	nameSize := len(SerializeDomainOrIP(answer.Name))
	return nameSize + 14 // 2+2+4+2+4 = 14 bytes plus the name
}
