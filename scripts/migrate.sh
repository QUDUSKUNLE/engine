#!/bin/bash

# Database migration script for Medicue
set -e

# Default values
DB_URL="${DB_URL:-postgres://medicue_user:medicue_password@localhost:5432/medicue?sslmode=disable}"
MIGRATION_PATH="${MIGRATION_PATH:-adapters/db/migrations}"
MIGRATE_BINARY="${MIGRATE_BINARY:-./bin/migrate}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[MIGRATE]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[MIGRATE]${NC} $1"
}

error() {
    echo -e "${RED}[MIGRATE]${NC} $1"
}

# Check if migrate binary exists
if [ ! -f "$MIGRATE_BINARY" ]; then
    error "Migration binary not found at: $MIGRATE_BINARY"
    error "Please run 'make setup' to install migration tools"
    exit 1
fi

# Check if migration files exist
if [ ! -d "$MIGRATION_PATH" ]; then
    error "Migration directory not found at: $MIGRATION_PATH"
    exit 1
fi

# Wait for database to be ready
wait_for_db() {
    log "Waiting for database to be ready..."
    
    # Extract database connection info
    DB_HOST=$(echo $DB_URL | sed 's|.*@\([^:]*\):.*|\1|')
    DB_PORT=$(echo $DB_URL | sed 's|.*:\([0-9]*\)/.*|\1|')
    DB_USER=$(echo $DB_URL | sed 's|.*://\([^:]*\):.*|\1|')
    
    # Wait for PostgreSQL to be ready
    for i in {1..30}; do
        if pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; then
            log "Database is ready!"
            return 0
        fi
        warn "Database not ready, waiting... ($i/30)"
        sleep 2
    done
    
    error "Database did not become ready in time"
    return 1
}

# Run migrations
run_migrations() {
    log "Running database migrations..."
    log "Database Host: ${DB_HOST%\?*}?***" # Hide query params for security
    log "Migration path: $MIGRATION_PATH"
    
    # Run migrations
    if $MIGRATE_BINARY -path="$MIGRATION_PATH" -database "$DB_URL" up; then
        log "✅ Migrations completed successfully!"
        return 0
    else
        error "❌ Migration failed!"
        return 1
    fi
}

# Get migration version
get_version() {
    log "Getting current migration version..."
    if $MIGRATE_BINARY -path="$MIGRATION_PATH" -database "$DB_URL" version; then
        return 0
    else
        warn "Could not get migration version (database might be empty)"
        return 1
    fi
}

# Force migration to specific version
force_version() {
    local version=$1
    if [ -z "$version" ]; then
        error "Please specify a version number"
        exit 1
    fi
    
    warn "Forcing migration to version: $version"
    if $MIGRATE_BINARY -path="$MIGRATION_PATH" -database "$DB_URL" force "$version"; then
        log "✅ Forced migration to version $version"
    else
        error "❌ Failed to force migration to version $version"
        exit 1
    fi
}

# Rollback migrations
rollback() {
    local steps=${1:-1}
    warn "Rolling back $steps migration(s)..."
    
    if $MIGRATE_BINARY -path="$MIGRATION_PATH" -database "$DB_URL" down "$steps"; then
        log "✅ Rollback completed successfully!"
    else
        error "❌ Rollback failed!"
        exit 1
    fi
}

# Main script logic
case "${1:-up}" in
    "up")
        wait_for_db
        run_migrations
        ;;
    "down")
        rollback "${2:-1}"
        ;;
    "version")
        get_version
        ;;
    "force")
        force_version "$2"
        ;;
    "wait")
        wait_for_db
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command] [options]"
        echo ""
        echo "Commands:"
        echo "  up                Run all pending migrations (default)"
        echo "  down [N]          Rollback N migrations (default: 1)"
        echo "  version           Show current migration version"
        echo "  force [VERSION]   Force migration to specific version"
        echo "  wait              Wait for database to be ready"
        echo "  help              Show this help message"
        echo ""
        echo "Environment variables:"
        echo "  DB_URL            Database connection string"
        echo "  MIGRATION_PATH    Path to migration files"
        echo "  MIGRATE_BINARY    Path to migrate binary"
        ;;
    *)
        error "Unknown command: $1"
        error "Use '$0 help' for usage information"
        exit 1
        ;;
esac
