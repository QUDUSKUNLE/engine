GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

.PHONY: setup generate migrate-up migrate-up-online migrate-down 

SQLC_VERSION=v1.27.0
MIGRATE_VERSION=v4.17.0

setup:
	go mod download
	GOBIN=$(GOBIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	GOBIN=$(GOBIN) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

generate:
	$(GOBIN)/sqlc generate -f adapters/db/sqlc.json

migrate-up:
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medivue?sslmode=disable" up

migrate-up-online:
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgresql://diagnoxix_postgres_user:arrxOcvRhZjKBNmQn9HVuDBNRErF3gT7@dpg-d2dmhj95pdvs73f2ce90-a.frankfurt-postgres.render.com/diagnoxix_postgres?sslmode=require" up

migrate-down:
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medivue?sslmode=disable" down

migrate-down-online:
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgresql://diagnoxix_postgres_user:arrxOcvRhZjKBNmQn9HVuDBNRErF3gT7@dpg-d2dmhj95pdvs73f2ce90-a.frankfurt-postgres.render.com/diagnoxix_postgres?sslmode=require" down

force-migrate:
	@if [ -z "$(V)" ]; then \
		echo "error: please specify version argument V"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medivue?sslmode=disable" force $(V)

migrate-rollback:
	@if [ -z "$(V)" ]; then \
		echo "error: please specify version argument V"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medivue?sslmode=disable" goto $(V)

migration-version:
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medivue?sslmode=disable" version

# Docker-based migration commands
migrate-docker-up:
	docker-compose -f docker-compose.migrate.yml up migrate

migrate-docker-down:
	docker-compose -f docker-compose.migrate.yml run --rm migrate ./bin/migrate -path=adapters/db/migrations -database "$(DB_URL)" down $(V)

migrate-docker-force:
	@if [ -z "$(V)" ]; then \
		echo "error: please specify version argument V"; \
		exit 1; \
	fi; \
	docker-compose -f docker-compose.migrate.yml run --rm migrate ./bin/migrate -path=adapters/db/migrations -database "$(DB_URL)" force $(V)

migrate-docker-version:
	docker-compose -f docker-compose.migrate.yml run --rm migrate ./bin/migrate -path=adapters/db/migrations -database "$(DB_URL)" version

# Production migration commands (use with caution)
migrate-prod-up:
	@if [ -z "$(PROD_DB_URL)" ]; then \
		echo "error: please set PROD_DB_URL environment variable"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "$(PROD_DB_URL)" up

migrate-prod-version:
	@if [ -z "$(PROD_DB_URL)" ]; then \
		echo "error: please set PROD_DB_URL environment variable"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "$(PROD_DB_URL)" version

# Create new migration
create-migration:
	@if [ -z "$(NAME)" ]; then \
		echo "error: please specify migration name with NAME argument"; \
		echo "example: make create-migration NAME=add_user_table"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate create -ext sql -dir adapters/db/migrations -seq $(NAME)

# Migration help
migrate-help:
	@echo "Available migration commands:"
	@echo "  migrate-up           - Run all pending migrations (local)"
	@echo "  migrate-down         - Rollback one migration (local)"
	@echo "  migrate-version      - Show current migration version (local)"
	@echo "  migrate-docker-up    - Run migrations using Docker"
	@echo "  migrate-docker-down  - Rollback migrations using Docker (V=number)"
	@echo "  migrate-docker-force - Force migration version using Docker (V=version)"
	@echo "  migrate-prod-up      - Run migrations on production (PROD_DB_URL required)"
	@echo "  create-migration     - Create new migration (NAME required)"
	@echo "  migrate-help         - Show this help"
