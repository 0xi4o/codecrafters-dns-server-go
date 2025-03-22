package main

import (
	"strings"
)

func SerializeDomainOrIP(input string) []byte {
	buf := []byte{}
	for _, label := range strings.Split(input, ".") {
		buf = append(buf, byte(len(label)))
		buf = append(buf, []byte(label)...)
	}
	buf = append(buf, byte(0))

	return buf
}

func DeserializeDomainOrIP(buf []byte) string {
	var labels []string
	pos := 0

	for pos < len(buf) {
		if pos >= len(buf) {
			break
		}
		labelLength := int(buf[pos])
		pos++

		if labelLength == 0 {
			break
		}
		label := string(buf[pos : pos+labelLength])
		labels = append(labels, label)

		pos += labelLength
	}

	return strings.Join(labels, ".")
}
