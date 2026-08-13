APP_NAME := production-go-api

DB_CONTAINER := production-go-postgres

DB_HOST := localhost
DB_PORT := 5433
DB_USER := production
DB_PASSWORD := production
DB_NAME := production_api

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# -----------------------
# Docker Database
# -----------------------

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-stop
docker-stop:
	docker compose stop

.PHONY: docker-restart
docker-restart:
	docker compose restart

.PHONY: docker-logs
docker-logs:
	docker compose logs -f

.PHONY: docker-ps
docker-ps:
	docker compose ps

# ===========
# Db shell
# =========

.PHONY: db-shell
db-shell:
	docker exec -it $(DB_CONTAINER) psql \
	-U $(DB_USER) \
	-d $(DB_NAME)

.PHONY: db-tables
db-tables:
	docker exec -it $(DB_CONTAINER) psql \
	-U $(DB_USER) \
	-d $(DB_NAME) \
	-c "\dt"

.PHONY: db-drop-all
db-drop-all:
	docker exec -it $(DB_CONTAINER) psql \
	-U $(DB_USER) \
	-d $(DB_NAME) \
	-c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

.PHONY: db-reset
db-reset:
	docker compose down -v
	docker compose up -d


# -----------------------
# Goose Migration
# -----------------------

.PHONY: migrate-create
migrate-create:
	goose -dir migrations create $(name) sql

.PHONY: migrate-up
migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

.PHONY: migrate-status
migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

.PHONY: migrate-reset
migrate-reset:
	goose -dir migrations postgres "$(DB_URL)" reset


# =========================
# Redis
# =========================

redis-up:
	docker compose up -d redis

redis-down:
	docker compose stop redis

redis-restart:
	docker compose restart redis

redis-logs:
	docker compose logs -f redis

redis-status:
	docker compose ps redis

redis-cli:
	docker exec -it production-redis redis-cli

redis-ping:
	docker exec production-redis redis-cli PING

redis-flush:
	docker exec production-redis redis-cli FLUSHDB

redis-keys:
	docker exec production-redis redis-cli KEYS '*'

redis-info:
	docker exec production-redis redis-cli INFO


# -----------------------
# Application
# -----------------------

.PHONY: run
run:
	go run ./cmd/api

.PHONY: build
build:
	go build -o bin/api ./cmd/api

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...