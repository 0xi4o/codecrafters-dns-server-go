package main

import (
	"fmt"
	"net"
)

func main() {
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

	fmt.Println("Listening for incoming UDP packets on port 2053...")

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := buf[:size]
		receivedHeader := DeserializeHeader(receivedData[:12])
		fmt.Printf("Received %d bytes from %s: %v\n", size, source, receivedData)

		header := NewHeader(receivedHeader.ID, receivedHeader.OPCODE, receivedHeader.RD)
		question := NewQuestion("codecrafters.io", 1, 1)
		answer := NewAnswer("codecrafters.io", 1, 1, 60, 4, "1.1.1.1")
		message := NewMessage(header, question, answer)
		response := message.Serialize()

		fmt.Printf("Sending response: %v\n", response)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
