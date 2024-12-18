# DNS-сервер

Этот проект представляет собой DNS-сервер, разработанный на языке программирования Go, который обеспечивает быструю и безопасную обработку DNS-запросов. Сервер предназначен для использования в средах, где требуется высокая производительность, надежность и защита от различных типов атак.

## Основные функции

### 1. Поддерживаемые функции
- **Обработка стандартных типов DNS-запросов**: Сервер поддерживает основные типы DNS-запросов.
- **Логирование запросов**: Вся активность сервера (входящие запросы, ответы, ошибки, состояние лимиттера и т.д.) логируется с возможностью настройки уровня детализации. Логи записываются в syslog.
- **Проксирование запросов**: Сервер может проксировать запросы на внешние DNS-серверы, обеспечивая гибкость и возможность использования дополнительных источников данных.

### 2. Особенности
- **Хранение данных в локальном файле**: Для хранения данных используется база данных SQLite, что обеспечивает простоту и удобство управления.
- **Встроенный кэш в памяти**: Для ускорения обработки запросов сервер использует встроенный кэш в памяти, что позволяет значительно сократить время ответа на повторяющиеся запросы.
- **Асинхронная обработка запросов**: Запросы обрабатываются асинхронно, что обеспечивает высокую производительность и отзывчивость сервера.

### 3. Механизмы защиты
- **ACL (список контроля доступа)**: Сервер поддерживает настройку ACL для ограничения доступа к определенным ресурсам.
- **Лимитирование числа запросов**: Для предотвращения DDoS-атак сервер ограничивает количество запросов от одного клиента в единицу времени.
- **Гибкая настройка лимиттера**: Параметры лимиттера (интервал времени, количество запросов и т.д.) могут быть настроены в соответствии с конкретными требованиями.
- **Защита от DNS amplification**: Сервер включает механизмы защиты от атак типа "DNS amplification", что обеспечивает безопасность от этого вида угроз.
- **Валидация входных данных**: Все входные данные проходят строгую валидацию, что предотвращает возможные инъекции и другие уязвимости.

## Установка и запуск

### 1. Запуск сервера
Для запуска DNS-сервера выполните следующую команду:
```bash
go run cmd/dns_server/main.go
```

### 2. Заполнение таблицы
После запуска сервера заполните таблицу данными с помощью скрипта:
```bash
python3 insert_records.py
```

### 3. Проверка работоспособности
Для проверки работоспособности сервера выполните запрос с помощью клиента:
```bash
python3 dns_client.py
```

## TODO

- Написать тесты.
- Добавить поддержку DNSSEC.
- Добавить возможность логирования во внешние хранилища (такие как Apache Kafka, Clickhouse или Sentry).
