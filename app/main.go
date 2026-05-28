package main

import (
	"encoding"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type DNSMarshelerUnmarshaler interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	resolver := flag.String("resolver", "1.1.1.1:53", "Address to forward the DNS query to")
	flag.Parse()

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

	var resolverUDPConn *net.UDPConn
	if *resolver != "" {
		parts := strings.Split(*resolver, ":")
		fmt.Println(parts)
		hostStr, portStr := parts[0], parts[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println("Invalid resolver port: ", err)
			return
		}
		resolverAddr := net.UDPAddr{
			IP:   net.ParseIP(hostStr),
			Port: port,
		}
		fmt.Println(resolverAddr)
		resolverUDPConn, err = net.DialUDP("udp", nil, &resolverAddr)
		if err != nil {
			fmt.Println("Cannot forward query to target: ", err)
			return
		}
	}

	buf := make([]byte, 512)
	resolverBuf := make([]byte, 512)

	for {
		_, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading data:", err)
			break
		}

		message := NewDNSMessage()
		err = message.UnmarshalBinary(buf)
		if err != nil {
			fmt.Println("Error unmarshaling message:", err)
			break
		}

		allAnswers := []DNSResourceRecord{}
		for _, question := range message.Questions {
			singleMessage := NewDNSMessage()

			singleMessage.Header.ID = message.Header.ID
			singleMessage.Header.RD = message.Header.RD
			singleMessage.Header.OPCODE = message.Header.OPCODE
			singleMessage.Header.QDCOUNT = 1

			singleMessage.Questions = []DNSQuestion{question}

			queryBuf, err := singleMessage.MarshalBinary()
			if err != nil {
				fmt.Println("Error marshaling single message:", err)
				break
			}

			_, err = resolverUDPConn.Write(queryBuf)
			if err != nil {
				fmt.Println("Error writing data to resolver:", err)
				break
			}

			rn, err := resolverUDPConn.Read(resolverBuf)
			if err != nil {
				fmt.Println("Error reading data from resolver:", err)
				break
			}

			resolverMessage := NewDNSMessage()
			err = resolverMessage.UnmarshalBinary(resolverBuf[:rn])
			if err != nil {
				fmt.Println("Error unmarshaling message from resolver:", err)
				break
			}

			allAnswers = append(allAnswers, resolverMessage.Answers...)
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
		// replace message.Answers with answers from resolver
		message.Answers = allAnswers
		message.Header.ANCOUNT = uint16(len(allAnswers))

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
