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

		message := NewDNSMessage()
		err = message.UnmarshalBinary(buf)
		if err != nil {
			fmt.Println("Error unmarshaling message:", err)
			break
		}

		message.Header.QR = true
		message.Header.AA = false
		message.Header.TC = false
		message.Header.RA = false
		message.Header.Z = 0
		if message.Header.OPCODE == 0 {
			message.Header.RCODE = 0
		} else {
			message.Header.RCODE = 4
		}
		message.Header.ANCOUNT = message.Header.QDCOUNT

		response, err := message.MarshalBinary()
		if err != nil {
			fmt.Println("Error unmarshaling message:", err)
			break
		}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
