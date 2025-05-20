SCHEMA_DIR="db/schema"
DB_DRIVER="postgres"
DB_DSN="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"

.PHONY: create-migration migrate-up migrate-down status

create-migration:
	goose -dir $(SCHEMA_DIR) create $(name) sql

migrate-up:
	goose -dir $(SCHEMA_DIR) $(DB_DRIVER) $(DB_DSN) up

migrate-down:
	goose -dir $(SCHEMA_DIR) $(DB_DRIVER) $(DB_DSN) down

status:
	goose -dir $(SCHEMA_DIR) $(DB_DRIVER) $(DB_DSN) status