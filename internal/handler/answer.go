package handler

import (
	"encoding/binary"
	"net"
	"strings"
)

// структура ответа пакета
type Answer struct {
	Name   string // доменное имя, к которому относится эта запись ресурса
	Type   Type   // содержит один из кодов типа RR
	Class  Class  // определяет класс данных в поле RDATA
	TTL    uint32 // время жизни записи
	Length uint32 // длинна поля RDATA
	Data   []byte // сообщение
}

// кодировка ответа
func (a Answer) Encode() []byte {
	var rrBytes []byte

	domain := a.Name
	parts := strings.Split(domain, ".")

	for _, label := range parts {
		if len(label) > 0 {
			rrBytes = append(rrBytes, byte(len(label)))
			rrBytes = append(rrBytes, []byte(label)...)
		}
	}
	rrBytes = append(rrBytes, 0x00)

	rrBytes = append(rrBytes, intToBytes(uint16(a.Type))...)
	rrBytes = append(rrBytes, intToBytes(uint16(a.Class))...)

	time := make([]byte, 4)
	binary.BigEndian.PutUint32(time, a.TTL)

	rrBytes = append(rrBytes, time...)
	rrBytes = append(rrBytes, intToBytes(a.Length)...)

	ipBytes, err := net.IPv4(a.Data[0], a.Data[1], a.Data[2], a.Data[3]).MarshalText()
	if err != nil {
		return nil
	}

	rrBytes = append(rrBytes, ipBytes...)

	return rrBytes
}
