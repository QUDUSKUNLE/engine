#!/bin/bash

# Entrypoint script for Medicue Docker container
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[ENTRYPOINT]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[ENTRYPOINT]${NC} $1"
}

error() {
    echo -e "${RED}[ENTRYPOINT]${NC} $1"
}

# Check if we should run migrations
RUN_MIGRATIONS=${RUN_MIGRATIONS:-true}

log "Starting Medicue application..."

# Run migrations if enabled
if [ "$RUN_MIGRATIONS" = "true" ]; then
    log "Running database migrations..."
    
    # Wait for database to be ready
    log "Waiting for database connection..."
    
    # Simple database connection check using migrate
    for i in $(seq 1 30); do
        if ./bin/migrate -path=adapters/db/migrations -database "$DATABASE_URL" version >/dev/null 2>&1; then
            log "Database connection established!"
            break
        elif [ $i -eq 30 ]; then
            error "Could not connect to database after 30 attempts"
            exit 1
        else
            warn "Database not ready, waiting... ($i/30)"
            sleep 2
        fi
    done
    
    # Run migrations
    log "Executing database migrations..."
    if ./bin/migrate -path=adapters/db/migrations -database "$DATABASE_URL" up; then
        log "✅ Database migrations completed successfully!"
    else
        error "❌ Database migration failed!"
        exit 1
    fi
else
    log "Skipping database migrations (RUN_MIGRATIONS=false)"
fi

# Start the application
log "Starting Medicue API server..."
exec "$@"
