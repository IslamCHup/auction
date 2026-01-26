<div align="center">

# Auction Platform (Gateway + Auction + User/Wallet + Notifications)

Микросервисная платформа аукционов с API Gateway, аукционом, кошельком пользователей и сервисом уведомлений.

[Демо: http://localhost:8080](http://localhost:8080)

</div>

---

## Коротко о проекте

- Платформа аукционов с торгами, пользовательскими кошельками и системой уведомлений на события.
- Архитектура: несколько Go-сервисов за API Gateway; взаимодействие по HTTP, события через Kafka; общая БД — PostgreSQL; на входе — JWT-аутентификация, rate limit и таймауты.

---

## Стек

- Язык/фреймворки: Go, Gin, GORM
- Инфраструктура: Docker, Docker Compose
- Хранилище: PostgreSQL
- Сообщения: Apache Kafka, Kafka UI
- Аутентификация: JWT 

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