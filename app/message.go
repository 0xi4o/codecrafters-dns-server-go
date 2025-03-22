package main

import (
	"bytes"
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

func NewMessage(header *Header, question *Question, answer *Answer) *Message {
	return &Message{
		Header:   *header,
		Question: *question,
		Answer:   *answer,
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

func NewHeader(id uint16, opcode uint8, rd uint8) *Header {
	header := &Header{
		ID:      id,
		QR:      1,
		OPCODE:  opcode,
		AA:      0,
		TC:      0,
		RD:      rd,
		RA:      0,
		Z:       0,
		RCODE:   0,
		QDCOUNT: 1,
		ANCOUNT: 1,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
	if opcode == 0 {
		header.RCODE = 0
	} else {
		header.RCODE = 4
	}

	return header
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

func DeserializeHeader(buf []byte) Header {
	flags := binary.BigEndian.Uint16(buf[2:4])
	return Header{
		ID:      binary.BigEndian.Uint16(buf[:2]),
		QR:      uint8(flags >> 15 & 0x1),
		OPCODE:  uint8(flags >> 11 & 0xF),
		AA:      uint8(flags >> 10 & 0x1),
		TC:      uint8(flags >> 9 & 0x1),
		RD:      uint8(flags >> 8 & 0x1),
		RA:      uint8(flags >> 7 & 0x1),
		Z:       uint8(flags >> 4 & 0x1),
		RCODE:   uint8(flags & 0xF),
		QDCOUNT: binary.BigEndian.Uint16(buf[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(buf[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(buf[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(buf[10:12]),
	}
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

func DeserializeQuestion(buf []byte) Question {
	name, flags, _ := bytes.Cut(buf, []byte{0})
	return Question{
		Name:  DeserializeDomainOrIP(name),
		Type:  binary.BigEndian.Uint16(flags[:2]),
		Class: binary.BigEndian.Uint16(flags[2:]),
	}
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
