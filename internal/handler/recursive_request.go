package handler

import (
	"net"
)

// ответ на рекурсивный запрос
func (h *DNSHandler) sendRecursiveRequest(question Question, source *net.UDPAddr) {
	dnsServer := "8.8.8.8:53" // используем Google DNS

	// формируем пакет для запроса на другой сервер
	query := make([]byte, 0)
	header := Header{
		PacketID: generateTransactionID(),
		QR:       0,
		OPCode:   0,
		RD:       1,
		QDCount:  1,
	}
	query = append(query, header.Encode()...)
	query = append(query, question.Encode()...)

	// устанавливаем соединение
	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		h.logger.Error("Failed to connect to DNS server: %v", err)
		return
	}
	defer conn.Close()

	// отправляем пакет
	_, err = conn.Write(query)
	if err != nil {
		h.logger.Error("Failed to send DNS query: %v", err)
		return
	}

	// считываем ответ
	response := make([]byte, 512)
	n, err := conn.Read(response)
	if err != nil {
		h.logger.Error("Failed to read DNS response: %v", err)
		return
	}

	// отправляем ответ клиенту
	_, err = h.listener.WriteToUDP(response[:n], source)
	if err != nil {
		h.logger.Error("Failed to send DNS response to client: %v", err)
	}
}
