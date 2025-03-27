package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	port := flag.String("port", "2053", "Port")
	resolver := flag.String("resolver", "", "Specify an address (ip:port) to forward the dns request")
	flag.Parse()

	udpResolverAddr, err := net.ResolveUDPAddr("udp", *resolver)
	if err != nil {
		fmt.Println("Failed to resolve resolver UDP address:", err)
		return
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%s", *port))
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

	fmt.Printf("Listening for incoming UDP packets on port %s...\n", *port)
	if *resolver != "" {
		fmt.Printf("Forwarding DNS queries to %v\n", *udpResolverAddr)
	} else {
		fmt.Println("No resolver found. Running in local-only mode...")
	}

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := buf[:size]
		receivedHeader := DeserializeHeader(receivedData[:12])
		offset := 12 // Start after the header
		header := NewHeader(receivedHeader.ID, receivedHeader.QR, receivedHeader.OPCODE, receivedHeader.RD, receivedHeader.QDCOUNT, 0)
		questions := []Question{}
		answers := []Answer{}
		for i := uint16(0); i < receivedHeader.QDCOUNT; i++ {
			receivedQuestion, newOffset := DeserializeQuestion(receivedData[offset:])
			question := NewQuestion(receivedQuestion.Name, receivedQuestion.Type, receivedQuestion.Class)
			questions = append(questions, *question)
			offset += newOffset
			if *resolver != "" {
				fmt.Printf("Forwarding A record query for: %s\n", receivedQuestion.Name)
				questionBuf := question.SerializeQuestion()
				message := NewMessage(header, []Question{receivedQuestion}, []Answer{})
				forward := message.Serialize()
				_, err = udpConn.WriteToUDP(forward, udpResolverAddr)
				if err != nil {
					fmt.Println("failed to forward message: ", err)
				}
				resolverBuf := make([]byte, 512)
				_, _, err := udpConn.ReadFromUDP(resolverBuf)
				if err != nil {
					fmt.Println("Error receiving data:", err)
					break
				}
				offset := 12 + len(questionBuf)
				answer := DeserializeAnswer(resolverBuf, offset)
				answers = append(answers, answer)
			} else {
				answer := NewAnswer(questions[i].Name, questions[i].Type, questions[i].Class, 60, 4, "1.1.1.1")
				answers = append(answers, *answer)
			}
		}

		header.QR = 1
		header.ANCOUNT = uint16(len(answers))
		message := NewMessage(header, questions, answers)
		response := message.Serialize()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
