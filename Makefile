ifneq (,$(wildcard .env))
include .env
export
endif

MIGRATE ?= migrate

.PHONY: run build vet db-create migrate-up migrate-down migrate-version migrate-create

run:
	go run .

build:
	go build ./...

vet:
	go vet ./...

db-create:
	@PGPASSWORD="$(POSTGRES_PASSWORD)" psql \
		--host="$(POSTGRES_HOST)" \
		--port="$(POSTGRES_PORT)" \
		--username="$(POSTGRES_USER)" \
		--dbname=postgres \
		--tuples-only \
		--command="SELECT 1 FROM pg_database WHERE datname = '$(POSTGRES_DB)'" | grep -q 1 || \
	PGPASSWORD="$(POSTGRES_PASSWORD)" createdb \
		--host="$(POSTGRES_HOST)" \
		--port="$(POSTGRES_PORT)" \
		--username="$(POSTGRES_USER)" \
		"$(POSTGRES_DB)"

migrate-up:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" down 1

migrate-version:
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" version

migrate-create:
	@test -n "$(name)" || (echo "usage: make migrate-create name=create_example" && exit 1)
	$(MIGRATE) create -ext sql -dir migrations -seq "$(name)"
