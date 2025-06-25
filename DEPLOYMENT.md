# Medicue Deployment Guide

This guide covers various deployment options for the Medicue API application.

## Prerequisites

- Docker and Docker Compose installed
- Git repository access
- Environment variables configured

## 🏠 Local Deployment

### Quick Start
```bash
./deploy-local.sh
```

### Manual Steps
```bash
# Build and run with Docker Compose
docker-compose up --build -d

# Check health
curl http://localhost:8080/v1/health

# View logs
docker-compose logs -f app
```

**Access Points:**
- API: http://localhost:8080
- Health Check: http://localhost:8080/v1/health
- Swagger Docs: http://localhost:8080/swagger/index.html
- Metrics: http://localhost:8080/metrics

## ☁️ Cloud Deployment Options

### 1. Railway (Recommended for MVP)

1. **Connect Repository**
   ```bash
   # Install Railway CLI
   npm install -g @railway/cli
   
   # Login and deploy
   railway login
   railway link
   railway up
   ```

2. **Environment Variables**
   Set these in Railway dashboard:
   - Copy values from `.env.production`
   - Update with your production credentials

3. **Database**
   - Railway will automatically provision PostgreSQL
   - Update `DB_URL` in environment variables

### 2. Render

1. **Connect Repository**
   - Go to [Render Dashboard](https://dashboard.render.com)
   - Connect your GitHub repository
   - Render will use `render.yaml` configuration

2. **Manual Setup**
   - Create Web Service from Docker
   - Set build command: `docker build -t medicue .`
   - Set start command: `./main`

### 3. DigitalOcean App Platform

1. **Create App**
   ```bash
   # Using doctl CLI
   doctl apps create --spec digitalocean-app.yaml
   ```

2. **Or use the web interface**
   - Connect GitHub repository
   - Choose Docker deployment
   - Configure environment variables

### 4. AWS ECS/Fargate

1. **Build and Push Image**
   ```bash
   # Build for AWS
   docker build -t medicue-api .
   
   # Tag and push to ECR
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ECR_URI
   docker tag medicue-api:latest YOUR_ECR_URI/medicue-api:latest
   docker push YOUR_ECR_URI/medicue-api:latest
   ```

2. **Create ECS Service**
   - Use the provided task definition
   - Configure load balancer
   - Set environment variables

### 5. Google Cloud Run

```bash
# Deploy to Cloud Run
gcloud run deploy medicue-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

### 6. Fly.io (Recommended for Production)

1. **Install Fly CLI**
   ```bash
   brew install flyctl
   ```

2. **Login to Fly.io**
   ```bash
   flyctl auth login
   ```

3. **Automated Deployment**
   ```bash
   ./deploy-fly.sh
   ```

4. **Manual Setup**
   ```bash
   # Create app and database
   flyctl launch --no-deploy --name medicue
   flyctl postgres create --name medicue-db
   flyctl postgres attach --app medicue medicue-db
   
   # Set environment variables
   flyctl secrets set JWT_SECRET_KEY="$(openssl rand -base64 32)" --app medicue
   flyctl secrets set ALLOW_ORIGINS="https://medicue.fly.dev" --app medicue
   
   # Deploy
   flyctl deploy --app medicue
   ```

5. **Monitor Deployment**
   ```bash
   flyctl logs --app medicue
   flyctl status --app medicue
   ```

### 7. VPS Deployment

```bash
# On your VPS
git clone https://github.com/your-username/medicue.git
cd medicue

# Copy production environment
cp .env.production .env

# Deploy
docker-compose up -d

# Setup reverse proxy (nginx example)
sudo nano /etc/nginx/sites-available/medicue
```

## 🔧 Configuration

### Environment Variables

**Required for Production:**
```env
PORT=8080
DB_URL=your-database-url
JWT_SECRET_KEY=your-secure-secret-key
ALLOW_ORIGINS=https://yourdomain.com
```

**Optional Services:**
- Cloudinary (file uploads)
- SendGrid (email)
- Paystack (payments)
- Google OAuth

### Database Setup

**Migrations are automatically run on container startup** by default.

**Manual Migration Commands:**
```bash
# Local development
make migrate-up              # Run all pending migrations
make migrate-down            # Rollback one migration
make migration-version       # Check current version
make create-migration NAME=add_something  # Create new migration

# Docker-based migrations
make migrate-docker-up       # Run migrations via Docker
make migrate-docker-version  # Check version via Docker

# Production migrations (manual)
PROD_DB_URL="your-prod-url" make migrate-prod-up

# Using migration script directly
./scripts/migrate.sh up      # Run migrations
./scripts/migrate.sh version  # Check version
./scripts/migrate.sh down 2   # Rollback 2 migrations
```

**Disable Auto-Migration:**
Set `RUN_MIGRATIONS=false` in your environment to skip automatic migrations on startup.

### SSL/TLS Certificate

**With Nginx:**
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d yourdomain.com
```

## 📊 Monitoring

### Health Checks
- Endpoint: `/v1/health`
- Returns service status and version

### Metrics
- Prometheus metrics: `/metrics`
- Grafana dashboard configuration included

### Logging
- Structured JSON logs
- Configurable log levels
- File and stdout output

## 🔒 Security Checklist

- [ ] Update JWT secret key
- [ ] Configure CORS origins
- [ ] Set secure database credentials
- [ ] Enable HTTPS
- [ ] Update API keys (Cloudinary, SendGrid, etc.)
- [ ] Configure rate limiting
- [ ] Set up firewall rules

## 🚀 CI/CD Pipeline

GitHub Actions workflow automatically:
1. Runs tests
2. Builds Docker image
3. Pushes to container registry
4. Deploys to configured platform

**To enable:**
1. Set repository secrets
2. Uncomment desired deployment section in `.github/workflows/deploy.yml`

## 🆘 Troubleshooting

### Common Issues

**Container won't start:**
```bash
# Check logs
docker-compose logs app

# Check environment variables
docker-compose exec app env
```

**Database connection failed:**
```bash
# Test database connectivity
docker-compose exec app pg_isready -h db -U medicue_user
```

**Health check fails:**
```bash
# Test locally
curl -v http://localhost:8080/v1/health

# Check application logs
docker-compose logs app
```

### Performance Tuning

**Database:**
- Connection pooling configured
- Indexes on frequently queried columns
- Regular VACUUM and ANALYZE

**Application:**
- Rate limiting enabled
- Request body size limits
- Graceful shutdown handling

## 📞 Support

For deployment issues:
1. Check logs first
2. Verify environment variables
3. Test database connectivity
4. Review security settings

**Useful Commands:**
```bash
# View all logs
docker-compose logs

# Restart specific service
docker-compose restart app

# Update and redeploy
git pull && docker-compose up --build -d

# Database backup
docker-compose exec db pg_dump -U medicue_user medicue > backup.sql
```
