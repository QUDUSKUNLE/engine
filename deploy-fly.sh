#!/bin/bash

# Fly.io deployment script for Medicue
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[FLY-DEPLOY]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[FLY-DEPLOY]${NC} $1"
}

error() {
    echo -e "${RED}[FLY-DEPLOY]${NC} $1"
}

info() {
    echo -e "${BLUE}[FLY-DEPLOY]${NC} $1"
}

# Check if flyctl is installed
check_flyctl() {
    if ! command -v flyctl &> /dev/null; then
        error "flyctl is not installed. Please install it first:"
        echo "brew install flyctl"
        exit 1
    fi
    log "✅ flyctl is installed"
}

# Check if user is logged in
check_auth() {
    if ! flyctl auth whoami &> /dev/null; then
        error "You are not logged in to Fly.io"
        info "Please run: flyctl auth login"
        exit 1
    fi
    log "✅ Authenticated with Fly.io"
}

# Create or update app
setup_app() {
    log "Setting up Fly.io application..."
    
    # Check if app already exists
    if flyctl apps list | grep -q "medicue"; then
        warn "App 'medicue' already exists. Updating configuration..."
    else
        log "Creating new Fly.io app..."
        if ! flyctl launch --no-deploy --copy-config --name medicue; then
            error "Failed to create Fly.io app"
            exit 1
        fi
    fi
}

# Set up PostgreSQL database
setup_database() {
    log "Setting up PostgreSQL database..."
    
    # Check if database already exists
    if flyctl postgres list | grep -q "medicue-db"; then
        warn "Database 'medicue-db' already exists"
        DB_URL=$(flyctl postgres connect -a medicue-db --command "echo \$DATABASE_URL" 2>/dev/null || echo "")
    else
        log "Creating PostgreSQL database..."
        if flyctl postgres create --name medicue-db --region ord; then
            log "✅ Database created successfully"
            # Attach database to app
            flyctl postgres attach --app medicue medicue-db
        else
            error "Failed to create database"
            exit 1
        fi
    fi
    
    log "✅ Database setup complete"
}

# Set environment variables
set_environment() {
    log "Setting environment variables..."
    
    # Required environment variables
    flyctl secrets set \
        JWT_SECRET_KEY="$(openssl rand -base64 32)" \
        JWT_EXPIRATION_HOURS="24" \
        ALLOW_ORIGINS="https://medicue.fly.dev" \
        RUN_MIGRATIONS="true" \
        --app medicue
    
    warn "⚠️  Please set the following secrets manually with your production values:"
    echo "flyctl secrets set CLOUDINARY_CLOUD_NAME=your-cloud-name --app medicue"
    echo "flyctl secrets set CLOUDINARY_API_KEY=your-api-key --app medicue"
    echo "flyctl secrets set CLOUDINARY_API_SECRET=your-api-secret --app medicue"
    echo "flyctl secrets set SENDGRID_API_KEY=your-sendgrid-key --app medicue"
    echo "flyctl secrets set PAYSTACK_SECRET_KEY=your-paystack-secret --app medicue"
    echo "flyctl secrets set PAYSTACK_PUBLIC_KEY=your-paystack-public --app medicue"
    echo "flyctl secrets set GOOGLE_CLIENT_ID=your-google-client-id --app medicue"
    echo "flyctl secrets set GOOGLE_CLIENT_SECRET=your-google-client-secret --app medicue"
    
    log "✅ Basic environment variables set"
}

# Deploy application
deploy_app() {
    log "Deploying application to Fly.io..."
    
    if flyctl deploy --app medicue; then
        log "✅ Deployment successful!"
        
        # Get app URL
        APP_URL=$(flyctl info --app medicue | grep "Hostname" | awk '{print $2}')
        
        info "🎉 Medicue is now deployed!"
        info "🌐 App URL: https://$APP_URL"
        info "🏥 Health Check: https://$APP_URL/v1/health"
        info "📚 API Docs: https://$APP_URL/swagger/index.html"
        info "📊 Metrics: https://$APP_URL/metrics"
        
        # Test health endpoint
        log "Testing health endpoint..."
        sleep 10
        if curl -f "https://$APP_URL/v1/health" > /dev/null 2>&1; then
            log "✅ Health check passed!"
        else
            warn "⚠️  Health check failed. Check logs with: flyctl logs --app medicue"
        fi
        
    else
        error "❌ Deployment failed!"
        error "Check logs with: flyctl logs --app medicue"
        exit 1
    fi
}

# Show useful commands
show_commands() {
    info "📝 Useful Fly.io commands:"
    echo ""
    echo "# Monitor application"
    echo "flyctl logs --app medicue"
    echo "flyctl status --app medicue"
    echo ""
    echo "# Database operations"
    echo "flyctl postgres connect --app medicue-db"
    echo "flyctl postgres list"
    echo ""
    echo "# Scaling"
    echo "flyctl scale count 2 --app medicue"
    echo "flyctl scale memory 2048 --app medicue"
    echo ""
    echo "# Environment variables"
    echo "flyctl secrets list --app medicue"
    echo "flyctl secrets set KEY=value --app medicue"
    echo ""
    echo "# Deployment"
    echo "flyctl deploy --app medicue"
    echo "flyctl releases --app medicue"
}

# Main deployment flow
main() {
    log "🚀 Starting Fly.io deployment for Medicue..."
    
    check_flyctl
    check_auth
    setup_app
    setup_database
    set_environment
    deploy_app
    show_commands
    
    log "🎉 Deployment process completed!"
}

# Handle command line arguments
case "${1:-deploy}" in
    "deploy")
        main
        ;;
    "database")
        check_flyctl
        check_auth
        setup_database
        ;;
    "env")
        check_flyctl
        check_auth
        set_environment
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  deploy      Full deployment process (default)"
        echo "  database    Set up database only"
        echo "  env         Set environment variables only"
        echo "  help        Show this help message"
        ;;
    *)
        error "Unknown command: $1"
        error "Use '$0 help' for usage information"
        exit 1
        ;;
esac
