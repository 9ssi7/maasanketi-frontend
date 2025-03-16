.PHONY: backend
backend:
	@$(MAKE) -C modules/backend up

.PHONY: web
web:
	@$(MAKE) -C modules/web up

.PHONY: up
up:
	docker compose up -d --build

.PHONY: down
down:
	docker compose down


# Database migration commands
.PHONY: migrate-up migrate-down migrate-create migrate-force migrate-version

# Path to migration files
MIGRATIONS_PATH=modules/backend/resources/migrations

# Database connection string (can be overridden via environment variable)
DB_URL ?= $(DB_CONN_STR)
# Default connection string for local development
ifeq ($(DB_URL),)
	DB_URL = postgres://postgres:s1cr1t@localhost:5432/postgres?sslmode=disable
endif

# Check if DB_URL is set
check-db-url:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL or DB_CONN_STR environment variable is not set"; \
		echo "Usage: DB_CONN_STR=postgres://username:password@localhost:5432/database_name make migrate-up"; \
		exit 1; \
	fi

# Apply all migrations
migrate-up: check-db-url
	@echo "Applying all migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# Rollback all migrations
migrate-down: check-db-url
	@echo "Rolling back all migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

# Rollback a specific number of migrations
migrate-down-steps: check-db-url
	@echo "Rolling back $(STEPS) migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down $(STEPS)

# Apply a specific number of migrations
migrate-up-steps: check-db-url
	@echo "Applying $(STEPS) migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up $(STEPS)

# Create a new migration
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required"; \
		echo "Usage: make migrate-create NAME=create_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration $(NAME)..."
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

# Force migration version
migrate-force: check-db-url
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required"; \
		echo "Usage: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@echo "Forcing migration version to $(VERSION)..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION)

# Show current migration version
migrate-version: check-db-url
	@echo "Current migration version:"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

# Install golang-migrate tool
.PHONY: install-migrate
install-migrate:
	@echo "Installing golang-migrate CLI tool..."
	@if command -v go >/dev/null 2>&1; then \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		echo "golang-migrate installed successfully. Make sure your Go bin directory is in your PATH."; \
	else \
		echo "Error: Go is not installed or not in your PATH."; \
		exit 1; \
	fi

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make up                    - Run the API server"
	@echo "  make down                  - Stop the API server"
	@echo ""
	@echo "Setup commands:"
	@echo "  make install-migrate       - Install golang-migrate CLI tool"
	@echo ""
	@echo "Migration commands:"
	@echo "  make migrate-up            - Apply all migrations"
	@echo "  make migrate-down          - Rollback all migrations"
	@echo "  make migrate-up-steps      - Apply a specific number of migrations (STEPS=n)"
	@echo "  make migrate-down-steps    - Rollback a specific number of migrations (STEPS=n)"
	@echo "  make migrate-create NAME=x - Create a new migration"
	@echo "  make migrate-force VERSION=x - Force migration version"
	@echo "  make migrate-version       - Show current migration version"
	@echo ""
	@echo "Environment variables:"
	@echo "  DB_CONN_STR - Database connection string (optional, defaults to local PostgreSQL)"
	@echo "  Example: DB_CONN_STR=postgres://username:password@localhost:5432/database_name make migrate-up"
