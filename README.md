# Auction Platform (Gateway + Auction + User/Wallet + Notifications)

Платформа аукционов, реализованная на микросервисной архитектуре.

---

## Коротко о проекте

Платформа аукционов, реализованная на микросервисной архитектуре.

Проект представляет собой микросервисную платформу аукциона. 
Система состоит из API Gateway и независимых сервисов (пользователи и кошелёк, аукционы, уведомления), развёрнутых в Docker-окружении.

---

## Стек

[![Go](https://img.shields.io/badge/Go-1.20%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?logo=go&logoColor=white)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-ORM-CE2D2D?logo=go&logoColor=white)](https://gorm.io/)
[![Docker](https://img.shields.io/badge/Docker-Engine-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![Docker Compose](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-DB-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-Event%20Streaming-231F20?logo=apachekafka&logoColor=white)](https://kafka.apache.org/)
[![Kafka UI](https://img.shields.io/badge/Kafka%20UI-Tool-231F20?logo=apachekafka&logoColor=white)](https://github.com/provectus/kafka-ui)
[![JWT](https://img.shields.io/badge/JWT-Auth-000000?logo=jsonwebtokens&logoColor=white)](https://jwt.io/)
[![slog](https://img.shields.io/badge/slog-Logging-4B5563?logo=go&logoColor=white)](https://pkg.go.dev/log/slog)
[![Makefile](https://img.shields.io/badge/Make-Build%20Tool-A3C51C?logo=gnu&logoColor=white)](https://www.gnu.org/software/make/)

---

## Схема 

![Архитектурная схема](./scheme.png)

---

## Быстрый старт

1) Создайте .env в корне проекта и в корне user-wallet-service с секретом JWT (нужен для Gateway и user-wallet):

```bash
cat > .env << 'EOF'
JWT_SECRET=change-me-super-secret
EOF

docker compose up -d --build
```

---


## Мой вклад в проект

### API Gateway
- Reverse Proxy для всех микросервисов (Auth/User, Wallet, Auction, Notification)
- JWT-валидация (HS256) с прокидыванием `X-User-Id` и `X-User-Role` через headers
- Rate limiting per-user (общий и отдельный для ставок)
- Таймауты на upstream-запросы
- Kafka consumers для событий `bid_placed` и `lot_completed`
- Асинхронное создание уведомлений для пользователей
- Хранение уведомлений в PostgreSQL
- Синхронные вызовы через Gateway с единым API
- Graceful shutdown Kafka consumers
- Структурированное логирование (`slog`)
- Создание Docker Compose


## Контакты

-  Email: [islam.ch@mail.ru](mailto:islam.ch@mail.ru)
- GitHub: [@IslamCHup](https://github.com/IslamCHup)
<br><br>

**Команда**
<br>

[Адам](https://github.com/warkaz16)<br>
[Ади Умаров](https://github.com/dasler-fw)

---
