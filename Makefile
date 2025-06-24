up:
	docker compose -f ./docker-compose.yml up --build -d

down:
	docker compose -f ./docker-compose.yml down

restart:
	down up

test:
	go test -v -cover ./...

swag:
	swag init -g cmd/server/main.go

.PHONY: up down restart test swag
