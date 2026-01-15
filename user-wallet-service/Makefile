.PHONY: run build test fmt vet lint tidy clean dev seed

GO           ?= go
BINARY       ?= app
CMD_MAIN     := ./cmd/app/main.go
CMD_SEED     := ./cmd/seed/main.go

run: ## Запуск основного приложения (HTTP-сервер)
	$(GO) run $(CMD_MAIN)

dev: ## Запуск в режиме разработки с hot reload (air)
	air -c .air.toml

build: ## Сборка бинарника приложения
	$(GO) build -o tmp/$(BINARY) $(CMD_MAIN)

test: ## Запуск всех тестов
	$(GO) test ./...

fmt: ## Форматирование кода
	$(GO) fmt ./...

vet: ## Статический анализ кода
	$(GO) vet ./...

lint: ## Линтинг кода с помощью golangci-lint
	golangci-lint run

tidy: ## Обновление зависимостей (go.mod / go.sum)
	$(GO) mod tidy

clean: ## Удаление собранных бинарников
	rm -rf tmp

seed: ## Заполнение БД начальными данными (пользователь и продукты)
	$(GO) run $(CMD_SEED)
