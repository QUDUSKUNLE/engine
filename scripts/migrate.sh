#!/bin/bash
set -e

# Default values
DATABASE_URL="${DATABASE_URL:-postgres://medicue_user:medicue_password@localhost:5432/medivue?sslmode=disable}"
MIGRATION_PATH="${MIGRATION_PATH:-adapters/db/migrations}"
MIGRATE_BINARY="${MIGRATE_BINARY:-./bin/migrate}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()   { echo -e "${GREEN}[MIGRATE]${NC} $1"; }
warn()  { echo -e "${YELLOW}[MIGRATE]${NC} $1"; }
error() { echo -e "${RED}[MIGRATE]${NC} $1"; }

wait_for_db() {
    log "Waiting for database to be ready..."
    DB_HOST=$(echo $DATABASE_URL | sed 's|.*@\([^:]*\):.*|\1|')
    DB_PORT=$(echo $DATABASE_URL | sed 's|.*:\([0-9]*\)/.*|\1|')
    DB_USER=$(echo $DATABASE_URL | sed 's|.*://\([^:]*\):.*|\1|')

    for i in {1..30}; do
        if pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" >/dev/null 2>&1; then
            log "Database is ready!"
            return 0
        fi
        warn "Database not ready, waiting... ($i/30)"
        sleep 2
    done

    error "Database did not become ready in time"
    return 1
}

run_migrations() {
    log "Running database migrations..."
    if $MIGRATE_BINARY -path="$MIGRATION_PATH" -database "$DATABASE_URL" up; then
        log "✅ Migrations completed successfully!"
    else
        error "❌ Migration failed!"
        exit 1
    fi
}

case "${1:-up}" in
    "up") wait_for_db && run_migrations ;;
    "down") $MIGRATE_BINARY -path="$MIGRATION_PATH" -database "$DATABASE_URL" down "${2:-1}" ;;
    "version") $MIGRATE_BINARY -path="$MIGRATION_PATH" -database "$DATABASE_URL" version ;;
    "force") $MIGRATE_BINARY -path="$MIGRATION_PATH" -database "$DATABASE_URL" force "$2" ;;
    "wait") wait_for_db ;;
    *) echo "Usage: $0 [up|down|version|force|wait]"; exit 1 ;;
esac
