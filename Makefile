GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

.PHONY: setup generate migrate-up migrate-down

SQLC_VERSION=v1.27.0
MIGRATE_VERSION=v4.17.0

setup:
	go mod download
	GOBIN=$(GOBIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	GOBIN=$(GOBIN) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

generate:
	$(GOBIN)/sqlc generate -f adapters/db/sqlc.json

migrate-up:
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medicue?sslmode=disable" up

migrate-down:
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medicue?sslmode=disable" down

force-migrate:
	@if [ -z "$(V)" ]; then \
		echo "error: please specify version argument V"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medicue?sslmode=disable" force $(V)

migrate-rollback:
	@if [ -z "$(V)" ]; then \
		echo "error: please specify version argument V"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medicue?sslmode=disable" goto $(V)

migration-version:
	$(GOBIN)/migrate -path="adapters/db/migrations" -database "postgres://abumuhsinah:abumuhsinah@localhost:5432/medicue?sslmode=disable" version
