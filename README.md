# Auction Platform (Gateway + Auction + User/Wallet + Notifications)

Микросервисная платформа аукционов с API Gateway, аукционом, кошельком пользователей и сервисом уведомлений.


---

## Коротко о проекте

- Платформа аукционов с торгами, пользовательскими кошельками и системой уведомлений на события.
- Архитектура: несколько Go-сервисов за API Gateway; взаимодействие по HTTP, события через Kafka; общая БД — PostgreSQL; на входе — JWT-аутентификация, rate limit и таймауты.

---

## Стек

[![Go](https://img.shields.io/badge/Go-1.20%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?logo=go&logoColor=white)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-ORM-CE2D2D?logo=go&logoColor=white)](https://gorm.io/)
[![Docker Compose](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-DB-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-Event%20Streaming-231F20?logo=apachekafka&logoColor=white)](https://kafka.apache.org/)
[![Kafka UI](https://img.shields.io/badge/Kafka%20UI-Tool-231F20?logo=apachekafka&logoColor=white)](https://github.com/provectus/kafka-ui)
[![JWT](https://img.shields.io/badge/JWT-Auth-000000?logo=jsonwebtokens&logoColor=white)](https://jwt.io/)

---

## Схема 

![Архитектурная схема](./scheme.png)

---

## Функционал

- Авторизация и профиль
	- Регистрация и вход (JWT)
	- Просмотр и изменение профиля: /api/users/me
- Каталог и лоты
	- Список лотов с пагинацией и фильтрами: статус, цена (min/max), даты окончания (min/max)
	- Просмотр лота по ID
	- Лоты конкретного пользователя
- Управление лотами (продавец)
	- Создание лота (draft), редактирование только в draft, публикация
- Участие в торгах
	- Создание ставки на активном лоте (валидация min шага и времени проведения)
	- Список ставок по лоту и мои ставки
- Кошелёк пользователя
	- Баланс и замороженные средства
	- Пополнение, заморозка/разморозка при перебитии ставки, списание у победителя
	- История транзакций
- Уведомления
	- Автоуведомления о перебитой ставке и завершении аукциона (через Kafka)
	- Список уведомлений, отметка как прочитано, счётчик непрочитанных
- Надёжность и защита
	- Rate limit на пользователя, таймауты запросов в Gateway, CORS
	- Периодическое завершение истёкших лотов (воркер) и служебные эндпойнты

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