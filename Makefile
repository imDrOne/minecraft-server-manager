export
	LOCAL_BIN:=$(CURDIR)/bin
	PATH:=$(LOCAL_BIN):$(PATH)

run:
	go mod tidy && go mod download && \
	go run ./cmd/app
.PHONY: run

# Prepare local environment
.PHONY: up-docker down-docker

up-docker:
	docker compose -p minecraft-server-manager up -d

down-docker:
	docker-compose stop

# Install external tools
.PHONY: bin-deps

bin-deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Migrations

.PHONY: migrate-create migrate-up

migrate-create:
	migrate create -ext sql -dir migrations "$(MIGRATE_NAME)"

migrate-up:
	go run ./cmd/migration

# Tools
.PHONY: sqlc-generate

sqlc-generate:
	sqlc generate -f sqlc.yaml

