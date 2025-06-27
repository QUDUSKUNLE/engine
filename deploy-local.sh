#!/bin/bash

# Local deployment script for Medicue
set -e

echo "🚀 Starting local deployment of Medicue..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Stop and remove existing containers
echo "🛑 Stopping existing containers..."
docker-compose down

# Remove old images
echo "🧹 Cleaning up old images..."
docker image prune -f

# Build and start services
echo "🔨 Building and starting services..."
docker-compose up --build -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are healthy
echo "🏥 Checking service health..."
if curl -f http://localhost:8080/v1/health > /dev/null 2>&1; then
    echo "✅ Medicue API is healthy and running at http://localhost:8080"
    echo "📚 API Documentation: http://localhost:8080/swagger/index.html"
    echo "📊 Metrics: http://localhost:8080/metrics"
    
    # Check migration status
    echo "🔍 Checking database migration status..."
    if docker-compose exec -T app ./bin/migrate -path=adapters/db/migrations -database "$DB_URL" version > /dev/null 2>&1; then
        MIGRATION_VERSION=$(docker-compose exec -T app ./bin/migrate -path=adapters/db/migrations -database "postgres://medicue_user:medicue_password@db:5432/medivue?sslmode=disable" version 2>/dev/null | tail -1)
        echo "📋 Current migration version: $MIGRATION_VERSION"
    else
        echo "⚠️  Could not check migration status"
    fi
else
    echo "❌ Service health check failed"
    echo "📋 Checking logs..."
    docker-compose logs app
    exit 1
fi

echo "🎉 Local deployment completed successfully!"
