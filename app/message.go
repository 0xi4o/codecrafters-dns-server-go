package main

import (
	"fmt"
)

type DNSMessage struct {
	Header     DNSHeader
	Questions  []DNSQuestion
	Answers    []DNSResourceRecord
	Authority  []DNSResourceRecord
	Additional []DNSResourceRecord
}

func NewDNSMessage() DNSMessage {
	header := DNSHeader{}
	questions := []DNSQuestion{}
	return DNSMessage{
		Header:    header,
		Questions: questions,
	}
}

func (m *DNSMessage) MarshalBinary() (data []byte, err error) {
	data = []byte{}
	responseHeader, err := m.Header.MarshalBinary()
	if err != nil {
		fmt.Println("Error marshaling header struct:", err)
	}
	responseQuestions := []byte{}
	responseAnswers := []byte{}
	for _, question := range m.Questions {
		responseQuestion, err := question.MarshalBinary()
		if err != nil {
			return []byte{}, fmt.Errorf("Error marshaling header struct: %w\n", err)
		}
		responseQuestions = append(responseQuestions, responseQuestion...)
		answer := DNSResourceRecord{
			Name:     question.Name,
			Type:     question.Type,
			Class:    question.Class,
			TTL:      60,
			RDLENGTH: 4,
			RDATA:    "8.8.8.8",
		}
		responseAnswer, err := answer.MarshalBinary()
		if err != nil {
			return []byte{}, fmt.Errorf("Error marshaling answer (resource record) struct: %w\n", err)
		}
		responseAnswers = append(responseAnswers, responseAnswer...)
	}
	data = append(data, responseHeader...)
	data = append(data, responseQuestions...)
	data = append(data, responseAnswers...)

	return data, nil
}

func (m *DNSMessage) UnmarshalBinary(buf []byte) error {
	err := m.Header.UnmarshalBinary(buf[:12])
	if err != nil {
		fmt.Println("Error unmarshaling header data:", err)
	}

	offset := 0
	for range m.Header.QDCOUNT {
		dnsQuestion := DNSQuestion{
			Offset: 12,
		}
		err = dnsQuestion.UnmarshalBinary(buf)
		if err != nil {
			fmt.Println("Error unmarshaling questions data:", err)
			break
		}
		offset += dnsQuestion.Offset
		m.Questions = append(m.Questions, dnsQuestion)
	}

	return nil
}
