package main

import (
	"dns-server/internal/acl"
	"dns-server/internal/config"
	"dns-server/internal/handler"
	"dns-server/internal/limiter"
	"dns-server/internal/logger"
	"dns-server/internal/migrator"
	"dns-server/internal/storage"
	"log"
	"time"
)

func main() {
	// загрузка конфига
	cfg, err := config.LoadConfig("configs", "config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// инициализация бд
	db, err := storage.NewSQLite("data/domains.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// запуск миграций
	err = migrator.RunMigrations(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	// инициализация acl менеджера и внесение адреса в список разрешённых
	aclManager := acl.NewACL()
	// добавляем разрешённые CIDR-диапазоны
	aclManager.AllowCIDR("192.168.1.0/24")
	aclManager.AllowCIDR("10.0.0.0/8")
	// добавляем разрешённые IP-адреса
	aclManager.AllowIP("127.0.0.1")
	// разрешаем рекурсивные запросы для определённых IP-адресов
	aclManager.AllowRecursion("127.0.0.1")

	// инициализация лимитера
	rateLimiter := limiter.NewRateLimiter(5, time.Minute)

	// инициализация логгера
	logger, err := logger.NewLogger(logger.InfoLvl)
	if err != nil {
		log.Fatal(err)
	}

	// инициализация и запуск сервера
	dnsHandler := handler.NewDNSHandler(cfg.Address, db, aclManager, rateLimiter, logger)
	dnsHandler.Start()
}
