package handler

import (
	"net"
)

func (h *DNSHandler) sendErrorResponse(header Header, source *net.UDPAddr, rCode uint16) {
	// кстановим параметры ответа
	header.QR = 1        // это ответ
	header.RCode = rCode // код ошибки
	header.ANCount = 0   // нет записей в ответе
	header.NSCount = 0
	header.ARCount = 0

	response := header.Encode() // кодируем заголовок

	// отправляем ответ клиенту
	_, err := h.listener.WriteToUDP(response, source)
	if err != nil {
		h.logger.Error("Failed to send error response: %v", err)
	}
}
