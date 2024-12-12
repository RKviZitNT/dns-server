package handler

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type Class uint16

const (
	_  Class = iota
	IN       // интернет
	CS       // класс CSNET
	CH       // класс CHAOS
	HS       // гесиод
)

type Type uint16

const (
	_     Type = iota
	A          // адрес хоста
	NS         // авторитетный сервер имен
	MD         // пункт назначения почты (устарело)
	MF         // сервер пересылки почты (устарело)
	CNAME      // каноническое имя псевдонима
	SOA        // обозначает начало зоны полномочий
	MB         // доменное имя почтового ящика (эксперементальный)
	MG         // член почтовой группы (эксперементальный)
	MR         // доменное имя переименования почты (эксперементальный)
	NULL       // нуль RR (эксперементальный)
	WKS        // известное описание службы
	PTR        // указатель доменного имени
	HINFO      // информация о хосте
	MINFO      // нформация о почтовом ящике или списке почты
	MX         // почтовый обмен
	TXT        // текстовые строки
	RRSIG      // подпись набора записей
)

// структара запроса (вопроса) пакета
type Question struct {
	QName  string // доменное имя, представленное в виде последовательности меток
	QType  Type   // двухоктетный код, который указывает тип запроса
	QClass Class  // двухоктетный код, который определяет класс запроса
}

// чтение запроса (вопроса) пакета
func ReadQuestion(buf []byte) Question {
	start := 0
	var nameParts []string

	for len := buf[start]; len != 0; len = buf[start] {
		start++
		nameParts = append(nameParts, string(buf[start:start+int(len)]))
		start += int(len)
	}
	questionName := strings.Join(nameParts, ".")
	start++

	questionType := binary.BigEndian.Uint16(buf[start : start+2])
	questionClass := binary.BigEndian.Uint16(buf[start+2 : start+4])

	q := Question{
		QName:  questionName,
		QType:  Type(questionType),
		QClass: Class(questionClass),
	}

	return q
}

// кодировка запроса (вопроса) для ответа
func (q Question) Encode() []byte {
	domain := q.QName
	parts := strings.Split(domain, ".")

	var buf bytes.Buffer

	for _, label := range parts {
		if len(label) > 0 {
			buf.WriteByte(byte(len(label)))
			buf.WriteString(label)
		}
	}
	buf.WriteByte(0x00)
	buf.Write(intToBytes(uint16(q.QType)))
	buf.Write(intToBytes(uint16(q.QClass)))

	return buf.Bytes()
}
