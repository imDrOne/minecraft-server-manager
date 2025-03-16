export
	LOCAL_BIN:=$(CURDIR)/bin
	PATH:=$(LOCAL_BIN):$(PATH)

run:
	go mod tidy && go mod download && \
	go run ./cmd/app

test:
	go test -v ./...

.PHONY: run test

# Prepare local environment
.PHONY: up-docker down-docker

up-docker:
	docker compose -p minecraft-server-manager up -d

down-docker:
	docker-compose stop

# Migrations

.PHONY: migrate-create migrate-up

migrate-create:
	go tool migrate create -ext sql -dir db/migrations "$(name)"

migrate-up:
	go run ./cmd/migrate

# Tools
.PHONY: sqlc-generate

sqlc-generate:
	go tool sqlc generate -f sqlc.yaml

