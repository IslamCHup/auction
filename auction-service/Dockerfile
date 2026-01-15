## Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Enable Go modules and turn off CGO for static binary
ENV CGO_ENABLED=0 GOOS=linux

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o auction-service ./cmd/app


## Stage 2: Run
FROM alpine:3.20

WORKDIR /app

# Create non-root user
RUN adduser -D appuser

# Copy binary from builder
COPY --from=builder /app/auction-service .

USER appuser

EXPOSE 8081

CMD ["./auction-service"]

