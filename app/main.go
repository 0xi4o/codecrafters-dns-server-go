package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
	port := flag.String("port", "2053", "Port")
	resolverAddr := flag.String("resolver", "", "Specify an address (ip:port) to forward the dns request")
	flag.Parse()

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
	if *resolverAddr != "" {
		fmt.Printf("Forwarding DNS queries to %s\n", *resolverAddr)
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

		if *resolverAddr != "" {
			mergedQuestions := []Question{}
			mergedAnswers := []Answer{}

			for _, question := range receivedQuestions {
				fmt.Printf("Forwarding A record query for: %s\n", question.Name)

				forwardHeader := NewHeader(receivedHeader.ID, receivedHeader.OPCODE, receivedHeader.RD, 1)
				forwardQuestion := NewQuestion(question.Name, question.Type, question.Class)
				forwardMessage := NewMessage(forwardHeader, []Question{*forwardQuestion}, []Answer{})

				response, err := forwardDNSQuery(*resolverAddr, forwardMessage.Serialize())
				if err != nil {
					fmt.Printf("Error forwarding DNS query: %v\n", err)
					continue
				}

				respHeader := DeserializeHeader(response[:12])

				mergedQuestions = append(mergedQuestions, *forwardQuestion)

				if respHeader.ANCOUNT > 0 {
					answerOffset := 12

					for i := uint16(0); i < respHeader.QDCOUNT; i++ {
						_, newOffset := DeserializeQuestion(response[answerOffset:])
						answerOffset += newOffset
					}

					for i := uint16(0); i < respHeader.ANCOUNT; i++ {
						answer, err := ExtractAnswerFromResponse(response, answerOffset)
						if err != nil {
							fmt.Printf("Error extracting answer: %v\n", err)
							break
						}
						mergedAnswers = append(mergedAnswers, answer)

						nameSize := len(SerializeDomainOrIP(answer.Name))
						answerOffset += nameSize + 14
					}
				}
			}

			if len(mergedQuestions) > 0 {
				mergedHeader := NewHeader(receivedHeader.ID, receivedHeader.OPCODE, receivedHeader.RD, uint16(len(mergedQuestions)))
				mergedHeader.ANCOUNT = uint16(len(mergedAnswers))

				mergedMessage := NewMessage(mergedHeader, mergedQuestions, mergedAnswers)
				mergedResponse := mergedMessage.Serialize()

				_, err = udpConn.WriteToUDP(mergedResponse, source)
				if err != nil {
					fmt.Println("Failed to send merged response:", err)
				}
			}
		} else {
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
}

func forwardDNSQuery(resolverAddr string, query []byte) ([]byte, error) {
	raddr, err := net.ResolveUDPAddr("udp", resolverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve resolver address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to resolver: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	_, err = conn.Write(query)
	if err != nil {
		return nil, fmt.Errorf("failed to send query to resolver: %v", err)
	}

	respBuf := make([]byte, 512)
	n, err := conn.Read(respBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to receive response from resolver: %v", err)
	}

	return respBuf[:n], nil
}
