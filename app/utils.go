package main

import "strings"

func SerializeDomainOrIP(input string) []byte {
	buf := []byte{}
	for _, label := range strings.Split(input, ".") {
		buf = append(buf, byte(len(label)))
		buf = append(buf, []byte(label)...)
	}
	buf = append(buf, byte(0))

	return buf
}
