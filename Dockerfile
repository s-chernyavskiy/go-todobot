FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todobot cmd/todobot/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

RUN adduser --disabled-password appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/todobot .

RUN mkdir -p /app/config

COPY config /app/config

ENV CONFIG_PATH="/app/config/local.yaml"

EXPOSE 8080

CMD ./todobot --config=$CONFIG_PATH
