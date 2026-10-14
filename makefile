include ./.env

REDIS_URL = redis://$(REDIS_USER):$(REDIS_PASS)@$(REDIS_HOST):$(REDIS_PORT)
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

MIGRATE_PATH=db/migrations
SEEDER_PATH=db/seeders/seeder.sql

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DB_URL) -path $(MIGRATE_PATH) up

migrate-down:
	migrate -database $(DB_URL) -path $(MIGRATE_PATH) down $(N)

db-seed:
	psql $(DB_URL) -f $(SEEDER_PATH)
	
db-reset:
	make migrate-down migrate-up db-seed flushdb

up:
	docker compose up -d --build

down:
	docker compose down --rmi local

ps:
	docker compose ps -a

logs:
	docker logs -f e-walet-backend