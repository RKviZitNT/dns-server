package handler

import (
	"bytes"
	"net"
	"time"
)

func (h *DNSHandler) sendAuthoritativeRequest(header Header, question Question, source *net.UDPAddr) {
	// проверяем кэш
	if cachedAddresses, found := h.cache.Get(question.QName); found {
		h.logger.Info("Cache hit for domain: %s", question.QName)
		h.sendResponseWithAddresses(header, question, source, cachedAddresses)
		return
	}

	// если в кэше нет, обращаемся к базе данных
	addresses, err := h.storage.GetAddress(question.QName)
	if err != nil {
		h.logger.Info("Domain not found: %s", question.QName)
		h.sendErrorResponse(header, source, 3) // код ошибки NXDomain
		return
	}

	// сохраняем результат в кэш с TTL (например, 5 минут)
	h.cache.Set(question.QName, addresses, 5*time.Minute)

	// отправляем ответ
	h.sendResponseWithAddresses(header, question, source, addresses)
}

func (h *DNSHandler) sendResponseWithAddresses(header Header, question Question, source *net.UDPAddr, addresses []string) {
	var answers []Answer
	for _, address := range addresses {
		ip := net.ParseIP(address).To4()
		if ip == nil {
			h.logger.Error("Invalid IP address: %s", address)
			h.sendErrorResponse(header, source, 2) // код ошибки ServFail
			return
		}

		answer := Answer{
			Name:   question.QName,
			Type:   A,
			Class:  IN,
			TTL:    300,
			Length: net.IPv4len,
			Data:   []byte{ip[0], ip[1], ip[2], ip[3]}, // преобразуем IP в массив байт
		}
		answers = append(answers, answer)
	}

	// обновляем заголовок
	header.QR = 1                         // ответ
	header.AA = 1                         // авторитетный ответ
	header.ANCount = uint16(len(answers)) // количество записей в ответе

	// формируем ответ
	var res bytes.Buffer
	res.Write(header.Encode())
	res.Write(question.Encode())
	for _, answer := range answers {
		res.Write(answer.Encode())
	}

	// отправляем ответ
	_, err := h.listener.WriteToUDP(res.Bytes(), source)
	if err != nil {
		h.logger.Error("Failed to send response: %v", err)
	}
}
