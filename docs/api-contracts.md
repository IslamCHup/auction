# API Contracts (short)

База:
- Формат: JSON, UTF-8
- Время: RFC3339 UTC (пример: 2026-01-26T12:00:00Z)
- Деньги: целые в мин. единицах (копейки)
- Версия: префикс /api

Базовые URL (локально):
- Gateway: http://localhost:8080
- Auction: http://localhost:8081
- User/Wallet: http://localhost:8082
- Notifications: http://localhost:8083

Маршрутизация (через Gateway):
- /api/auth/*, /api/users/*, /api/wallet/* → User/Wallet
- /api/lots*, /api/users/:id/{lots,bids} → Auction
- /api/notifications/* → Notifications

Аутентификация:
- Authorization: Bearer <JWT>
- Gateway валидирует JWT и пробрасывает X-User-Id

Ошибки (унифицированно):
- { "error": { "code", "message", "request_id" } }
- Коды: unauthorized, forbidden, not_found, validation_error, conflict, rate_limited, internal

Пагинация:
- Query: page (>=1), page_size (1..100)
- Ответ: { data: [...], pagination: { page, page_size, total, has_more } }


## 1 Auth / Users
- POST /api/auth/register → 201 { user, token } (409, если email занят)
- POST /api/auth/login → 200 { user, token }
- GET /api/users/me (JWT) → 200 User
- PATCH /api/users/me (JWT) → 200 User

User (основные поля): id, full_name, email, role, created_at


## 2 Wallet
- GET /api/wallet/me (JWT) → 200 { user_id, balance, frozen_balance, currency, updated_at }
- GET /api/wallet/transactions?type=&page=&page_size= (JWT) → 200 { data, pagination }
  - type: deposit|withdraw|freeze|unfreeze|charge|refund


## 3 Lots & Bids
Сущности (ключевые поля):
- Lot: id, title, description, start_price, current_price, min_step, status(draft|active|finished|canceled), start_at, end_at, seller_id, winner_id, current_bid_id, created_at, updated_at
- Bid: id, lot_id, user_id, amount, created_at

Эндпойнты:
- GET /api/lots?status=&price_min=&price_max=&end_at_from=&end_at_to=&seller_id=&search=&page=&page_size= → 200 { data, pagination }
- GET /api/lots/:id → 200 Lot | 404
- POST /api/lots (JWT) → 201 Lot (status=draft)
  - Валидации: end_at > start_at, min_step > 0, start_price > 0
- PATCH /api/lots/:id (JWT владелец, только draft) → 200 Lot | 409
- POST /api/lots/:id/publish (JWT владелец) → 200 Lot(status=active) | 409