package handler

import (
	"dns-server/internal/acl"
	"dns-server/internal/limiter"
	"dns-server/internal/logger"
	memorycache "dns-server/internal/memory_cache"
	"dns-server/internal/storage"
	"log"
	"net"
	"time"
)

type DNSHandler struct {
	address  string
	storage  *storage.SQLite
	listener *net.UDPConn
	acl      *acl.ACL
	limiter  *limiter.RateLimiter
	logger   *logger.Logger
	cache    *memorycache.MemoryCache
}

func NewDNSHandler(address string, db *storage.SQLite, acl *acl.ACL, limiter *limiter.RateLimiter, logger *logger.Logger) *DNSHandler {
	return &DNSHandler{
		address: address,
		storage: db,
		acl:     acl,
		limiter: limiter,
		logger:  logger,
		cache:   memorycache.NewMemoryCache(),
	}
}

// запуск сервера
func (h *DNSHandler) Start() {
	// cоздаем UDP-сокет
	udpAddr, err := net.ResolveUDPAddr("udp", h.address)
	if err != nil {
		log.Fatal("failed to resolve udp address", err)
	}

	h.listener, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("failed to bind to address", err)
	}
	defer h.listener.Close()

	log.Printf("Started server on %s", h.address)

	// запуск фоновой очистки кэша
	go h.startCacheCleanup()

	// обработка входящих запросов
	buffer := make([]byte, 512)
	for {
		size, source, err := h.listener.ReadFromUDP(buffer)
		if err != nil {
			h.logger.Error("Failed to read from UDP: %v", err)
			continue
		}

		// обрабатываем запрос в отдельной горутине
		go h.handleRequest(buffer[:size], source)
	}
}

// автоматический фоновый процесс очистки кэша
func (h *DNSHandler) startCacheCleanup() {
	ticker := time.NewTicker(1 * time.Minute) // интервал очистки
	defer ticker.Stop()

	// используем for range для обработки тиков
	for range ticker.C {
		h.cache.Cleanup()
	}
}
