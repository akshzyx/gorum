BACKEND_DIR=backend
FRONTEND_DIR=frontend

run-backend:
	cd $(BACKEND_DIR) && air

run-frontend:
	cd $(FRONTEND_DIR) && npm run dev

run:
	make -j2 run-backend run-frontend

install:
	cd $(BACKEND_DIR) && go mod tidy
	cd $(FRONTEND_DIR) && npm install