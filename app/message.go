package main

import (
	"encoding/binary"
)

type Message struct {
	Header   Header
	Question Question
	Answer   Answer
}

type Header struct {
	ID      uint16
	QR      uint8
	OPCODE  uint8
	AA      uint8
	TC      uint8
	RD      uint8
	RA      uint8
	Z       uint8
	RCODE   uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

type Answer struct {
	Name   string
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   string
}

func NewMessage() *Message {
	return &Message{
		Header:   *NewHeader(),
		Question: *NewQuestion("codecrafters.io", 1, 1),
		Answer:   *NewAnswer("codecrafters.io", 1, 1, 60, 4, "1.1.1.1"),
	}
}

func (m *Message) Serialize() []byte {
	messageBuf := []byte{}
	headerBuf := m.Header.SerializeHeader()
	questionBuf := m.Question.SerializeQuestion()
	answerBuf := m.Answer.SerializeAnswer()

	messageBuf = append(messageBuf, headerBuf...)
	messageBuf = append(messageBuf, questionBuf...)
	messageBuf = append(messageBuf, answerBuf...)

	return messageBuf
}

func NewHeader() *Header {
	return &Header{
		ID:      1234,
		QR:      1,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      0,
		RA:      0,
		Z:       0,
		RCODE:   0,
		QDCOUNT: 1,
		ANCOUNT: 1,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
}

func (header *Header) SerializeHeader() []byte {
	headerBuf := make([]byte, 0, 12)
	flags := uint16(0)
	flags |= uint16(header.QR) << 15
	flags |= uint16(header.OPCODE) << 11
	flags |= uint16(header.AA) << 10
	flags |= uint16(header.TC) << 9
	flags |= uint16(header.RD) << 8
	flags |= uint16(header.RA) << 7
	flags |= uint16(header.Z) << 4
	flags |= uint16(header.RCODE)
	headerBuf = binary.BigEndian.AppendUint16(headerBuf, header.ID)
	headerBuf = binary.BigEndian.AppendUint16(headerBuf, flags)
	headerBuf = binary.BigEndian.AppendUint16(headerBuf, header.QDCOUNT)
	headerBuf = binary.BigEndian.AppendUint16(headerBuf, header.ANCOUNT)
	headerBuf = binary.BigEndian.AppendUint16(headerBuf, header.NSCOUNT)
	headerBuf = binary.BigEndian.AppendUint16(headerBuf, header.ARCOUNT)

	return headerBuf
}

func NewQuestion(name string, qtype uint16, qclass uint16) *Question {
	return &Question{
		Name:  name,
		Type:  qtype,
		Class: qclass,
	}
}

func (question *Question) SerializeQuestion() []byte {
	questionBuf := []byte{}
	qNameBuf := SerializeDomainOrIP(question.Name)
	questionBuf = append(questionBuf, qNameBuf...)
	questionBuf = binary.BigEndian.AppendUint16(questionBuf, question.Type)
	questionBuf = binary.BigEndian.AppendUint16(questionBuf, question.Class)

	return questionBuf
}

func NewAnswer(name string, atype, aclass uint16, ttl uint32, length uint16, data string) *Answer {
	return &Answer{
		Name:   name,
		Type:   atype,
		Class:  aclass,
		TTL:    ttl,
		Length: length,
		Data:   data,
	}
}

func (answer *Answer) SerializeAnswer() []byte {
	answerBuf := []byte{}
	aNameBuf := SerializeDomainOrIP(answer.Name)
	answerBuf = append(answerBuf, aNameBuf...)
	answerBuf = binary.BigEndian.AppendUint16(answerBuf, answer.Type)
	answerBuf = binary.BigEndian.AppendUint16(answerBuf, answer.Class)
	answerBuf = binary.BigEndian.AppendUint32(answerBuf, answer.TTL)
	answerBuf = binary.BigEndian.AppendUint16(answerBuf, answer.Length)
	dataBuf := SerializeDomainOrIP(answer.Data)
	answerBuf = append(answerBuf, dataBuf...)

	return answerBuf
}
