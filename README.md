project-name/
├── docker-compose.yaml       # Вся инфраструктура
├── .gitignore
├── README.md
├── docs/
│   └── api-contracts.md      # Согласованные контракты
│
├── gateway/                  # API Gateway
│   ├── cmd/app/main.go
│   ├── internal/
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── user-service/             # User Service
│   ├── cmd/app/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── services/
│   │   └── transport/
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── order-service/            # Order Service
│   └── ... (аналогичная структура)
│
└── notification-service/     # Notification Service
    └── ... (аналогичная структура)