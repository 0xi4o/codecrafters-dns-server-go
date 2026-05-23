package main

import (
	"encoding"
	"fmt"
	"net"
)

type DNSMarshelerUnmarshaler interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type DNSMessage struct {
	Header     DNSHeader
	Question   DNSQuestion
	Answers    []_DNSResourceRecord
	Authority  []_DNSResourceRecord
	Additional []_DNSResourceRecord
}

type _DNSResourceRecord struct{}

// Ensures gofmt doesn't remove the "net" import in stage 1 (feel free to remove this!)
var _ = net.ListenUDP

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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

		var dnsHeader DNSHeader
		err = dnsHeader.UnmarshalBinary(buf[:12])
		if err != nil {
			fmt.Println("Error unmarshaling header data:", err)
			break
		}

		var dnsQuestion DNSQuestion
		err = dnsQuestion.UnmarshalBinary(buf[12:])
		if err != nil {
			fmt.Println("Error unmarshaling header data:", err)
			break
		}

		dnsHeader.QDCOUNT = 1
		responseHeader, err := dnsHeader.MarshalBinary()
		if err != nil {
			fmt.Println("Error marshaling header struct:", err)
			break
		}
		responseQuestion, err := dnsQuestion.MarshalBinary()
		if err != nil {
			fmt.Println("Error marshaling header struct:", err)
			break
		}

		// Create an empty response
		response := []byte{}
		response = append(response, responseHeader...)
		response = append(response, responseQuestion...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
