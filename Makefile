# Run backend (Air)
run-backend:
	cd backend && air

# Run frontend (Next.js)
run-frontend:
	cd frontend && npm run dev

# Run both concurrently
run:
	make -j2 run-backend run-frontend

# Install frontend deps
install-frontend:
	cd frontend && npm install

# Tidy Go modules
tidy-backend:
	cd backend && go mod tidy