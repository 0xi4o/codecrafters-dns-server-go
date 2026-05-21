package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type _DNSMessage struct {
	Header     DNSHeader
	Questions  []_DNSQuestion
	Answers    []_DNSResourceRecord
	Authority  []_DNSResourceRecord
	Additional []_DNSResourceRecord
}

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

type _DNSQuestion struct{}
type _DNSResourceRecord struct{}

func DecodeHeader(buf []byte) DNSHeader {
	// part1 - QR (1 bit), OPCODE (4 bits), AA (1 bit), TC (1 bit), RD (1 bit)
	// part2 - RA (1 bit), Z (3 bits), and RCODE (4 bits)
	part1, part2 := buf[2], buf[3]
	return DNSHeader{
		ID:      binary.BigEndian.Uint16(buf[:2]),
		QR:      part1&0b1000_0000 != 0,
		OPCODE:  (part1 >> 3) & 0x0F,
		AA:      part1&0b0000_0100 != 0,
		TC:      part1&0b0000_0010 != 0,
		RD:      part1&0b0000_0001 != 0,
		RA:      part2&0b1000_0000 != 0,
		Z:       (part2 >> 4) & 0x07,
		RCODE:   part2 & 0x0F,
		QDCOUNT: binary.BigEndian.Uint16(buf[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(buf[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(buf[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(buf[10:12]),
	}
}

func EncodeHeader(h DNSHeader) []byte {
	header := make([]byte, 12)
	fmt.Print(len(header))
	binary.BigEndian.PutUint16(header[:2], h.ID)
	var part1 byte
	var part2 byte
	// QR, OPCODE, AA, TC, and RD
	part1 |= 0b1000_0000
	// RA, Z, and RCODE
	part2 |= 0b1000_0000
	header[2] = part1
	header[3] = part2
	binary.BigEndian.PutUint16(header[4:6], h.QDCOUNT)
	binary.BigEndian.PutUint16(header[6:8], h.ANCOUNT)
	binary.BigEndian.PutUint16(header[8:10], h.NSCOUNT)
	binary.BigEndian.PutUint16(header[10:12], h.ARCOUNT)
	return header
}

// Ensures gofmt doesn't remove the "net" import in stage 1 (feel free to remove this!)
var _ = net.ListenUDP

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// TODO: Uncomment the code below to pass the first stage
	//
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		_, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		header := DecodeHeader(buf[:12])
		fmt.Printf("Header: %v\n", header)

		responseHeader := EncodeHeader(header)

		// Create an empty response
		response := []byte{}
		response = append(response, responseHeader...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
