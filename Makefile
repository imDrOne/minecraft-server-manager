export
	LOCAL_BIN:=$(CURDIR)/bin
	PATH:=$(LOCAL_BIN):$(PATH)

up-docker:
	docker compose -p minecraft-server-manager up -d
.PHONY: up-docker

down-docker:
	docker-compose stop
.PHONY: down-docker

bin-deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
.PHONY: bin-deps

migrate-create:
	migrate create -ext sql -dir db/migrations "$(MIGRATE_NAME)"
.PHONY: migrate-create

migrate-up:
	go run ./cmd/migration
.PHONY: migrate-up

sqlc-generate:
	sqlc generate -f ./db/sqlc.yaml
.PHONY: sqlc-generate
