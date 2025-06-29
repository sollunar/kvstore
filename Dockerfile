FROM golang:1.24.4 AS builder

WORKDIR /app

RUN apt-get update && \
    apt-get install -y libssl-dev pkg-config build-essential && \
    apt-get clean

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o kvstore-api ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y libssl3 && \
    apt-get clean

COPY --from=builder /app/kvstore-api .

CMD ["./kvstore-api"]

