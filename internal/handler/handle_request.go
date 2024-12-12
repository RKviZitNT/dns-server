package handler

import (
	"net"
)

// обработка типа запроса и формирование ответа
func (h *DNSHandler) handleRequest(data []byte, source *net.UDPAddr) {
	h.logger.Info("Incoming request from IP: %s", source.IP.String())

	// проверяем авторизацию
	if !h.acl.IsAllowed(source.IP.String()) {
		h.logger.Warning("Access denied for IP: %s", source.IP.String())
		header := ReadHeader(data[:12])
		h.sendErrorResponse(header, source, 5) // код ошибки Refused
		return
	}

	// проверяем лимиты
	if !h.limiter.Allow(source.IP.String()) {
		h.logger.Warning("Rate limit exceeded for IP: %s", source.IP.String())
		header := ReadHeader(data[:12])
		h.sendErrorResponse(header, source, 5) // код ошибки Refused
		return
	}

	if len(data) < 12 {
		h.logger.Error("Received malformed request")
		header := ReadHeader(data[:12])
		h.sendErrorResponse(header, source, 1) // код ошибки FormErr
		return
	}

	// считываем заголовок и запрос из пакета
	header := ReadHeader(data[:12])
	question := ReadQuestion(data[12:])

	if question.QName == "" || question.QType == 0 {
		h.logger.Error("Malformed question section in request")
		h.sendErrorResponse(header, source, 1) // код ошибки FormErr
		return
	}

	// проверяем, разрешён ли рекурсивный запрос для данного IP-адреса
	if header.RD == 1 && !h.acl.IsRecursionAllowed(source.IP.String()) {
		h.logger.Warning("Recursive query denied for IP: %s", source.IP.String())
		h.sendErrorResponse(header, source, 5) // код ошибки Refused
		return
	}

	// определяем тип запроса
	switch h.deternineQueryType(header, question) {
	case "authoritative":
		h.logger.Info("Handling authoritative request for domain: %s", question.QName)
		h.sendAuthoritativeRequest(header, question, source)
	case "redirect":
		h.logger.Info("Handling redirect request for domain: %s", question.QName)
		h.sendRedirectRequest(header, question, source)
	case "recursive":
		h.logger.Info("Handling recursive request for domain: %s", question.QName)
		h.sendRecursiveRequest(question, source)
	case "notFound":
		h.logger.Warning("Domain not found: %s", question.QName)
		h.sendErrorResponse(header, source, 3) // код ошибки NotFound
	default:
		h.logger.Error("Unknown request type for domain: %s", question.QName)
		h.sendErrorResponse(header, source, 4) // код ошибки NotImp
	}
}

// обработка типа запроса
func (h *DNSHandler) deternineQueryType(header Header, question Question) string {
	if h.cache.IsExists(question.QName) {
		return "authoritative"
	}
	if auth, redir := h.storage.DomainExists(question.QName); auth {
		if redir {
			return "redirect"
		}
		return "authoritative"
	}
	if header.RD == 1 {
		return "recursive"
	}
	return "notFound"
}
