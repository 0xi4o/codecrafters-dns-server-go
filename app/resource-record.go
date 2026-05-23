package main

import (
	"encoding/binary"
)

type DNSResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    string
}

func (rr *DNSResourceRecord) MarshalBinary() (data []byte, err error) {
	data = []byte{}
	encodedDomainName, err := encodeDomainName(rr.Name)
	if err != nil {
		return []byte{}, err
	}
	data = append(data, encodedDomainName...)
	data = append(data, 0x00)
	data = binary.BigEndian.AppendUint16(data, rr.Type)
	data = binary.BigEndian.AppendUint16(data, rr.Class)
	data = binary.BigEndian.AppendUint32(data, rr.TTL)
	data = binary.BigEndian.AppendUint16(data, rr.RDLENGTH)
	encodedARecord, err := encodeARecord(rr.RDATA)
	if err != nil {
		return []byte{}, err
	}
	data = append(data, encodedARecord...)
	return data, nil
}

// func (rr *DNSResourceRecord) UnmarshalBinary(buf []byte) error {
// 	name, offset, err := decodeDomainName(buf, 0)
// 	if err != nil {
// 		return err
// 	}

// 	rr.Name = name
// 	rr.Type = binary.BigEndian.Uint16(buf[offset : offset+2])
// 	rr.Class = binary.BigEndian.Uint16(buf[offset+2 : offset+4])
// 	rr.TTL = binary.BigEndian.Uint32(buf[offset+4 : offset+8])
// 	rr.RDLENGTH = binary.BigEndian.Uint16(buf[offset+8 : offset+10])
// 	return nil
// }
