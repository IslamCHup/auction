<div align="center">

# Auction Platform (Gateway + Auction + User/Wallet + Notifications)

Микросервисная платформа аукционов с API Gateway, аукционом, кошельком пользователей и сервисом уведомлений.

[Демо: http://localhost:8080](http://localhost:8080)

</div>

---

## Коротко о проекте

- Платформа аукционов с торгами в реальном времени, кошельками пользователей и уведомлениями через Kafka.
- Архитектура из нескольких Go-сервисов, объединённых через API Gateway и Docker Compose.

---

## Стек

- Язык/фреймворки: Go (Gin), GORM
- Инфраструктура: Docker, Docker Compose
- Хранилище: PostgreSQL
- Сообщения: Apache Kafka (sarama / kafka-go), Kafka UI
- Аутентификация: JWT (через Gateway и user-wallet-service)

---

## Сервисы и порты

- Gateway: http://localhost:8080
- Auction Service: http://localhost:8081
- User/Wallet Service: http://localhost:8082
- Notification Service: http://localhost:8083
- Kafka UI: http://localhost:8085

Контракты API: см. docs/api-contracts.md  
Postman-коллекция: postman/collections/Auction Gateway API.postman_collection.json

---

## Быстрый старт

Требования: Docker, Docker Compose, Make (опционально), Go ≥ 1.20 (для локального запуска без контейнеров).

1) Создайте .env в корне с секретом JWT (нужен для Gateway и user-wallet):

```bash
cat > .env << 'EOF'
JWT_SECRET=change-me-super-secret
EOF
```

2) Запустите всю инфраструктуру:

```bash
docker compose up -d --build
```

3) Проверьте доступность Gateway:

- Откройте http://localhost:8080
- Авторизационные эндпоинты идут через /api/auth/* (проксируются в user-wallet-service)

Полезно:

- Kafka UI: http://localhost:8085
- PostgreSQL: localhost:5432 (user=postgres password=12345 db=app_db)

Остановка:

```bash
docker compose down -v
```

---

## Примеры запуска по сервисам (локально)

- Auction Service:

```bash
cd auction-service
make run
```

- User/Wallet Service:

```bash
cd user-wallet-service
make run
```

- Notification Service:

```bash
cd notification-service
go run cmd/app/main.go
```

Примечание: для локального запуска без Docker убедитесь, что запущены PostgreSQL и Kafka, а переменные окружения DATABASE_URL, KAFKA_BROKERS, JWT_SECRET заданы корректно (см. docker-compose.yaml для примеров значений).

---

## Структура репозитория (сокращённо)

```
./
├── docker-compose.yaml
├── docs/
│   └── api-contracts.md
├── gateway/
├── auction-service/
├── user-wallet-service/
└── notification-service/
```


---

## Лицензия

MIT