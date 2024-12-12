package handler

import (
	"bytes"
	"net"
)

func (h *DNSHandler) sendRedirectRequest(header Header, question Question, source *net.UDPAddr) {
	// обращаемся к базе данных
	target, err := h.storage.GetRedirectAddress(question.QName)
	if err != nil {
		h.logger.Info("Domain not found or no redirect: %s", question.QName)
		h.sendErrorResponse(header, source, 3) // код ошибки NXDomain
		return
	}

	// преобразуем доменное имя в формат DNS
	cnameData := encodeDomainName(target)

	// формируем CNAME-ответ
	answer := Answer{
		Name:   question.QName,
		Type:   CNAME,
		Class:  IN,
		TTL:    300,
		Length: uint32(len(cnameData)),
		Data:   cnameData,
	}

	// обновляем заголовок
	header.QR = 1      // это ответ
	header.AA = 1      // авторитетный ответ
	header.ANCount = 1 // одна запись в ответе

	// формируем полный пакет ответа
	var res bytes.Buffer
	res.Write(header.Encode())   // кодируем заголовок
	res.Write(question.Encode()) // кодируем вопрос
	res.Write(answer.Encode())   // кодируем ответ

	// отправляем ответ клиенту
	_, err = h.listener.WriteToUDP(res.Bytes(), source)
	if err != nil {
		h.logger.Error("Failed to send redirect response: %v", err)
	}
}
