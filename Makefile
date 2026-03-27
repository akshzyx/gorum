BACKEND_DIR=backend
FRONTEND_DIR=frontend

# backend

run-backend:
	cd $(BACKEND_DIR) && air

migrate-up:
	cd $(BACKEND_DIR) && make migrate-up

migrate-down:
	cd $(BACKEND_DIR) && make migrate-down

migrate:
	cd $(BACKEND_DIR) && make migrate $(filter-out $@,$(MAKECMDGOALS))

reset:
	cd $(BACKEND_DIR) && make reset

seed:
	cd $(BACKEND_DIR) && make seed

dev:
	cd $(BACKEND_DIR) && make dev

sqlc:
	cd $(BACKEND_DIR) && sqlc generate

# frontend

run-frontend:
	cd $(FRONTEND_DIR) && npm run dev

# both

run:
	make -j2 run-backend run-frontend

install:
	cd $(BACKEND_DIR) && go mod tidy
	cd $(FRONTEND_DIR) && npm install

# Prevent Make from treating extra args as targets
%:
	@: