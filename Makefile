export
	LOCAL_BIN:=$(CURDIR)/bin
	PATH:=$(LOCAL_BIN):$(PATH)

.PHONY: run test unit_test integration_test coverage_report

run:
	go mod tidy && go mod download && \
	go run ./cmd/app

unit_test:
	go test -v ./internal/...

integration_test:
	go test -v ./test/integration...

test: unit_test integration_test

coverage_report:
	go test -p=1 -coverpkg=./... -count=1 -coverprofile=.coverage.out ./...
	go tool cover -html .coverage.out -o .coverage.html
	open ./.coverage.html

# Prepare local environment
.PHONY: up-docker down-docker

up-docker:
	docker compose -p minecraft-server-manager up -d

down-docker:
	docker-compose stop

# Migrations

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

