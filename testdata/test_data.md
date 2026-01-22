# Примеры данных для тестирования

## Регистрация пользователей

### Покупатель
```json
{
  "email": "buyer@example.com",
  "password": "password123",
  "role": "buyer"
}
```

### Продавец
```json
{
  "email": "seller@example.com",
  "password": "password123",
  "role": "seller"
}
```

## Login

```json
{
  "email": "buyer@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "buyer@example.com",
    "full_name": "",
    "role": "buyer"
  }
}
```

## Wallet Operations

### Deposit
```json
{
  "amount": 1000,
  "description": "Initial deposit"
}
```

**Headers:**
```
X-User-Id: 1
Content-Type: application/json
```

### Freeze Funds
```json
{
  "amount": 100,
  "description": "Bid placement for lot #1"
}
```

### Unfreeze Funds
```json
{
  "amount": 100,
  "description": "Bid retracted"
}
```

## Lots

### Create Lot
```json
{
  "title": "Vintage Watch",
  "description": "Beautiful vintage Swiss watch from 1980s",
  "starting_price": 500,
  "seller_id": 2,
  "ends_at": "2026-01-29T10:00:00Z"
}
```

### Update Lot
```json
{
  "title": "Updated Vintage Watch",
  "description": "Beautiful vintage Swiss watch from 1980s - UPDATED",
  "starting_price": 600
}
```

### List Lots with Filters
```
GET /api/lots?page=1&limit=10&status=active&min_price=100&max_price=1000
```

**Parameters:**
- `page` - Номер страницы (по умолчанию 1)
- `limit` - Количество на странице (по умолчанию 10, максимум 100)
- `status` - Статус: `draft`, `active`, `completed` (опционально)
- `min_price` - Минимальная цена (опционально)
- `max_price` - Максимальная цена (опционально)
- `min_end_date` - Минимальная дата окончания RFC3339 (опционально)
- `max_end_date` - Максимальная дата окончания RFC3339 (опционально)

## Bids

### Create Bid
```json
{
  "amount": 550,
  "bidder_id": 1
}
```

**Endpoint:** `POST /api/lots/1/bids`

### Get All Bids for Lot
```
GET /api/lots/1/bids
```

### Get All Bids by User
```
GET /api/users/1/bids
```

## Status Codes

| Code | Значение |
|------|----------|
| 200 | OK - Успешно |
| 201 | Created - Создано |
| 400 | Bad Request - Неверный запрос |
| 401 | Unauthorized - Не авторизован |
| 404 | Not Found - Не найдено |
| 500 | Internal Server Error - Ошибка сервера |

## Сценарии тестирования

### Сценарий 1: Полный цикл аукциона

1. **Регистрация пользователей**
   - Зарегистрировать покупателя
   - Зарегистрировать продавца

2. **Подготовка средств**
   - Пополнить кошелек покупателя (1000)
   - Заморозить средства (100)

3. **Создание лота**
   - Создать лот от продавца
   - Проверить, что лот в статусе `draft` или `active`

4. **Размещение ставок**
   - Создать ставку от покупателя (550)
   - Получить все ставки для лота
   - Проверить, что ставка видна

5. **Завершение**
   - Получить информацию о лоте
   - Проверить финальный статус

### Сценарий 2: Фильтрация лотов

1. Создать несколько лотов с разными ценами
2. Запросить лоты с фильтрами:
   - `status=active`
   - `min_price=100`
   - `max_price=1000`
3. Проверить корректность результатов

### Сценарий 3: Управление профилем

1. Зарегистрировать пользователя
2. Выполнить login и получить JWT token
3. Получить информацию о текущем пользователе (`GET /api/auth/me`)
4. Обновить профиль (`PUT /api/auth/me`)
5. Проверить обновленные данные

## Отладка

### Проверка логов контейнеров
```bash
docker compose logs -f gateway
docker compose logs -f auction
docker compose logs -f user-wallet
```

### Проверка статуса контейнеров
```bash
docker compose ps
```

### Очистка данных
```bash
docker compose down -v
```
