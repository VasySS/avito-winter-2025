ifneq (,$(wildcard ./.env))
    include .env
    export
endif

PG_IP := ${DATABASE_HOST}
PG_USER := ${DATABASE_USER}
PG_PASS := ${DATABASE_PASSWORD}
PG_DB := ${DATABASE_NAME}
PG_PORT := ${DATABASE_PORT}
PG_URL := postgres://${PG_USER}:${PG_PASS}@${PG_IP}:${PG_PORT}/${PG_DB}

MAIN_FILE := ./cmd/server/main.go

.PHONY: run
run:
	go run ${MAIN_FILE}

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: goose-install
goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: goose-add
goose-add:
	goose -dir ./migrations postgres "${PG_URL}" create hw-migration sql

.PHONY: goose-up
goose-up:
	goose -dir ./migrations postgres "${PG_URL}" up

.PHONY: goose-down
goose-down:
	goose -dir ./migrations postgres "${PG_URL}" down-to 0

.PHONY: goose-status
goose-status:
	goose -dir ./migrations postgres "${PG_URL}" status
