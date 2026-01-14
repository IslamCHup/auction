# build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# runtime stage
FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache add ca-certificates \
    && adduser -D appuser

USER appuser

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]