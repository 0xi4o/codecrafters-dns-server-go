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

	for _, question := range m.Questions {
		responseQuestion, err := question.MarshalBinary()
		if err != nil {
			return []byte{}, fmt.Errorf("error marshaling question: %w", err)
		}
		responseQuestions = append(responseQuestions, responseQuestion...)
	}

	responseAnswers := []byte{}
	for _, answer := range m.Answers {
		responseAnswer, err := answer.MarshalBinary()
		if err != nil {
			return []byte{}, fmt.Errorf("error marshaling answer: %w", err)
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

	offset := 12
	for range m.Header.QDCOUNT {
		dnsQuestion := DNSQuestion{
			Offset: offset,
		}
		err = dnsQuestion.UnmarshalBinary(buf)
		if err != nil {
			fmt.Println("Error unmarshaling questions data:", err)
			break
		}
		offset = dnsQuestion.Offset
		m.Questions = append(m.Questions, dnsQuestion)
	}

	answerOffset := offset
	for range m.Header.ANCOUNT {
		dnsAnswer := DNSResourceRecord{
			Offset: answerOffset,
		}
		err = dnsAnswer.UnmarshalBinary(buf)
		if err != nil {
			fmt.Println("Error unmarshaling questions data:", err)
			break
		}
		answerOffset = dnsAnswer.Offset
		m.Answers = append(m.Answers, dnsAnswer)
	}

	return nil
}
