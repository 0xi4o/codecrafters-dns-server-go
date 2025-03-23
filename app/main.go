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
		// fmt.Printf("Received data: %v\n", receivedData)
		receivedHeader := DeserializeHeader(receivedData[:12])
		// fmt.Printf("Received header: %v\n", receivedHeader)
		receivedQuestions := []Question{}
		offset := 12 // Start after the header
		for i := uint16(0); i < receivedHeader.QDCOUNT; i++ {
			receivedQuestion, newOffset := DeserializeQuestion(receivedData[offset:])
			receivedQuestions = append(receivedQuestions, receivedQuestion)
			offset += newOffset
		}

		// fmt.Printf("Received questions: %v\n", receivedQuestions)

		header := NewHeader(receivedHeader.ID, receivedHeader.OPCODE, receivedHeader.RD, receivedHeader.QDCOUNT)
		questions := []Question{}

		answers := []Answer{}
		for i := uint16(0); i < header.QDCOUNT; i++ {
			question := NewQuestion(receivedQuestions[i].Name, receivedQuestions[i].Type, receivedQuestions[i].Class)
			questions = append(questions, *question)
			answer := NewAnswer(receivedQuestions[i].Name, receivedQuestions[i].Type, receivedQuestions[i].Class, 60, 4, "1.1.1.1")
			answers = append(answers, *answer)
		}
		message := NewMessage(header, questions, answers)
		response := message.Serialize()

		// fmt.Printf("Sending header: %v\n", header)
		// fmt.Printf("Sending questions: %v\n", questions)
		// fmt.Printf("Sending answers: %v\n", answers)
		// fmt.Printf("Sending response: %v\n", response)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
