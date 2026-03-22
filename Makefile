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

# Prevent Make from thinking the second word is a separate target
%:
	@:
