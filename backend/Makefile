# Makefile for Gorum - Simple migration commands

include .env
export

MIGRATIONS_DIR=./internal/db/migrations

.DEFAULT_GOAL := help

help:
	@echo "Usage:"
	@echo "  make migrate-up                # Run all pending migrations"
	@echo "  make migrate-down              # Rollback last migration"
	@echo "  make migrate create_name       # Create a new migration"
	@echo "  make reset                     # Reset database (drop + migrate)"
	@echo "  make seed                      # Seed database with dummy data"
	@echo "  make dev                       # Reset + seed (full fresh setup)"

# Run migrations up
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

# Rollback last migration
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down

# Create new migration (usage: make migrate create_tweets)
migrate:
	@if [ "$(word 2,$(MAKECMDGOALS))" = "" ]; then \
		echo "Please provide migration name, e.g. make migrate create_tweets"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(word 2,$(MAKECMDGOALS))

# Reset database (drop schema + run migrations)
reset:
	./scripts/reset_db.sh

# Seed database (run Go seeder)
seed:
	go run scripts/seed.go

# Full dev setup (clean DB + seed data)
dev:
	$(MAKE) reset
	$(MAKE) seed

# Prevent Make from thinking the second word is a separate target
%:
	@: