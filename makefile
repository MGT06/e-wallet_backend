include ./.env

REDIS_URL = redis://$(REDIS_USER):$(REDIS_PASSWORD)@$(REDIS_HOST):$(REDIS_PORT)
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

MIGRATE_PATH=db/migrations
SEEDER_PATH=db/seeders

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DB_URL) -path $(MIGRATE_PATH) up

migrate-down:
	migrate -database $(DB_URL) -path $(MIGRATE_PATH) down $(N)

db-seed:
	cat $(SEEDER_PATH)/*.sql | psql $(DB_URL)


db-reset:
	make migrate-down migrate-up db-seed

up:
	docker compose -f deployment/docker-compose.yml up -d --build
	make logs

down:
	docker compose -f deployment/docker-compose.yml down --rmi local

reset:
	make down
	make up

ps:
	docker compose -f deployment/docker-compose.yml ps -a

logs:
	docker logs -f e-walet-backend

swag:
	swag init -g cmd/main.go -o docs