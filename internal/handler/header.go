package handler

import (
	"encoding/binary"
)

// структура заголовка пакета согласно RFC
type Header struct {
	PacketID uint16 // идентификатор
	QR       uint16 // является ли запросом
	OPCode   uint16 // тип запроса
	AA       uint16 // авторитетный ответ
	TC       uint16 // усечено ли сообщение из-за длины
	RD       uint16 // рекурсивный запрос
	RA       uint16 // доступна ли поддержка рекурсии
	Z        uint16 // поддержка DNSSEC (DO-бит)
	RCode    uint16 // код ошибки
	QDCount  uint16 // количество записей в разделе запросов (вопросов)
	ANCount  uint16 // количество записей ресурсов в разделе ответов
	NSCount  uint16 // количество записей ресурсов сервера имен в разделе записей полномочий
	ARCount  uint16 // количество записей ресурсов в разделе дополнительных записей
}

// чтение заголовка пакета
func ReadHeader(buffer []byte) Header {
	h := Header{
		PacketID: uint16(buffer[0])<<8 | uint16(buffer[1]),
		QR:       1,
		OPCode:   uint16((buffer[2] << 1) >> 4),
		AA:       uint16((buffer[2] << 5) >> 7),
		TC:       uint16((buffer[2] << 6) >> 7),
		RD:       uint16((buffer[2] << 7) >> 7),
		RA:       uint16(buffer[3] >> 7),
		Z:        uint16((buffer[3] << 1) >> 5),
		QDCount:  uint16(buffer[4])<<8 | uint16(buffer[5]),
		ANCount:  uint16(buffer[5])<<8 | uint16(buffer[7]),
		NSCount:  uint16(buffer[8])<<8 | uint16(buffer[9]),
		ARCount:  uint16(buffer[10])<<8 | uint16(buffer[11]),
	}

	// сервер может обрабатывать только стандартные значение, OPCode == 0 является стандартным запросом
	if h.OPCode == 0 {
		h.RCode = 0
	} else {
		h.RCode = 4
	}

	return h
}

// кодировка заголовка для ответа
func (h Header) Encode() []byte {
	dnsHeader := make([]byte, 12)

	var flags uint16 = 0
	flags = h.QR<<15 | h.OPCode<<11 | h.AA<<10 | h.TC<<9 | h.RD<<8 | h.RA<<7 | h.Z<<4 | h.RCode

	binary.BigEndian.PutUint16(dnsHeader[0:2], h.PacketID)
	binary.BigEndian.PutUint16(dnsHeader[2:4], flags)
	binary.BigEndian.PutUint16(dnsHeader[4:6], h.QDCount)
	binary.BigEndian.PutUint16(dnsHeader[6:8], h.ANCount)
	binary.BigEndian.PutUint16(dnsHeader[8:10], h.NSCount)
	binary.BigEndian.PutUint16(dnsHeader[10:12], h.ARCount)

	return dnsHeader
}
