package handler

import (
	"bytes"
	"encoding/binary"
	"strings"

	"golang.org/x/exp/rand"
)

// изменение типа int в uint16 и uint32
func intToBytes(n interface{}) []byte {
	switch t := n.(type) {
	case uint16:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, t)
		return buf
	case uint32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, t)
		return buf
	default:
		return nil
	}
}

// генератор ID
func generateTransactionID() uint16 {
	return uint16(rand.Intn(65536))
}

// кодируем имя домена
func encodeDomainName(domain string) []byte {
	parts := strings.Split(domain, ".")
	var encoded bytes.Buffer
	for _, part := range parts {
		encoded.WriteByte(byte(len(part))) // длина сегмента
		encoded.WriteString(part)          // сегмент имени
	}
	encoded.WriteByte(0) // завершающий нулевой байт
	return encoded.Bytes()
}
