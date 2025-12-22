
include .env
export $(shell sed 's/=.*//' .env)

SCHEMA_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

migration:
	migrate create -ext sql -dir ./pkg/database/migration $(table)

migrate-up:
	migrate -path ./pkg/database/migration -database "$(SCHEMA_URL)" up

migrate-down:
	echo y | migrate -path ./pkg/database/migration -database "$(SCHEMA_URL)" down

SEED_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)&x-migrations-table=seed_migrations

seed-up:
	migrate -path ./pkg/database/seeds -database "$(SEED_URL)" up

seed-down:
	echo y | migrate -path ./pkg/database/seeds -database "$(SEED_URL)" down

seeder:
	migrate create -ext sql -dir ./pkg/database/seeds -seq $(table)