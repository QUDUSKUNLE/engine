#!/bin/bash
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log()   { echo -e "${GREEN}[ENTRYPOINT]${NC} $1"; }
warn()  { echo -e "${YELLOW}[ENTRYPOINT]${NC} $1"; }
error() { echo -e "${RED}[ENTRYPOINT]${NC} $1"; }

RUN_MIGRATIONS=${RUN_MIGRATIONS:-true}

log "Starting Medicue application..."

if [ "$RUN_MIGRATIONS" = "true" ]; then
    log "Running database migrations..."
    ./migrate.sh up
else
    warn "Skipping database migrations (RUN_MIGRATIONS=false)"
fi

log "Launching API server..."
exec "$@"
