# AI Context Gap Tracker - Deployment Guide

## Overview

This guide provides step-by-step instructions for deploying the AI Context Gap Tracker system in various environments.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start with Docker](#quick-start-with-docker)
3. [Local Development Setup](#local-development-setup)
4. [Production Deployment](#production-deployment)
5. [Environment Configuration](#environment-configuration)
6. [Troubleshooting](#troubleshooting)
7. [Monitoring and Maintenance](#monitoring-and-maintenance)

## Prerequisites

### Required Software

- **Docker & Docker Compose**: Version 20.10+ and 1.29+
- **Go**: Version 1.21+ (for local development)
- **Python**: Version 3.11+ (for NLP service development)
- **Git**: For cloning the repository

### System Requirements

- **Memory**: Minimum 4GB RAM (8GB recommended)
- **Storage**: 10GB free space
- **Network**: Ports 8080, 5000, 5432, 6379 available

## Quick Start with Docker

### 1. Clone the Repository

```bash
git clone https://github.com/cliffordotieno/ai-context-gap-tracker.git
cd ai-context-gap-tracker
```

### 2. Start Services

```bash
# Build and start all services
make docker-run

# Or use Docker Compose directly
docker-compose up -d
```

### 3. Verify Installation

```bash
# Check service health
make health-check

# Run system tests
make system-test
```

### 4. Access Services

- **Main API**: http://localhost:8080
- **NLP Service**: http://localhost:5000
- **API Documentation**: http://localhost:8080/api/v1/health

## Local Development Setup

### 1. Install Dependencies

```bash
# Set up development environment
make dev-setup

# Or manually:
go mod tidy
cd python-nlp && pip install -r requirements.txt
```

### 2. Start Infrastructure Services

```bash
# Start only PostgreSQL and Redis
docker-compose up -d postgres redis
```

### 3. Run Services Locally

```bash
# Terminal 1: Start main Go service
make run

# Terminal 2: Start NLP service
cd python-nlp
uvicorn main:app --host 0.0.0.0 --port 5000 --reload
```

### 4. Run Tests

```bash
# Run Go tests
make test

# Run system integration tests
make system-test

# Run example usage
make example
```

## Production Deployment

### 1. Environment Preparation

```bash
# Create production environment file
cp .env.example .env.production

# Edit configuration
nano .env.production
```

### 2. Build Production Images

```bash
# Build optimized images
make prod-build
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build
```

### 3. Deploy with Docker Compose

```bash
# Deploy to production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 4. Set Up Reverse Proxy (Optional)

```nginx
# nginx.conf
upstream ai_context_tracker {
    server localhost:8080;
}

upstream nlp_service {
    server localhost:5000;
}

server {
    listen 80;
    server_name your-domain.com;

    location /api/ {
        proxy_pass http://ai_context_tracker;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /nlp/ {
        proxy_pass http://nlp_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Environment Configuration

### Required Environment Variables

```bash
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=ai_context_tracker
DB_USER=tracker_user
DB_PASSWORD=tracker_password
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

# Server Configuration
HTTP_PORT=8080
GRPC_PORT=9090

# NLP Service Configuration
NLP_SERVICE_URL=http://nlp-service:5000
NLP_TIMEOUT=30
```

### Optional Environment Variables

```bash
# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Security
JWT_SECRET=your-secret-key
API_KEY=your-api-key

# Performance
MAX_CONNECTIONS=100
TIMEOUT=30s
```

## Troubleshooting

### Common Issues

#### 1. Services Not Starting

```bash
# Check logs
docker-compose logs

# Check specific service
docker-compose logs ai-context-tracker

# Check system resources
docker stats
```

#### 2. Database Connection Issues

```bash
# Test database connection
docker-compose exec postgres psql -U tracker_user -d ai_context_tracker

# Check database logs
docker-compose logs postgres
```

#### 3. Redis Connection Issues

```bash
# Test Redis connection
docker-compose exec redis redis-cli ping

# Check Redis logs
docker-compose logs redis
```

#### 4. Port Conflicts

```bash
# Check port usage
sudo netstat -tulpn | grep :8080

# Stop conflicting services
sudo systemctl stop service-name
```

### Performance Issues

#### 1. Memory Usage

```bash
# Monitor memory usage
docker stats

# Adjust memory limits in docker-compose.yml
services:
  ai-context-tracker:
    mem_limit: 1g
```

#### 2. CPU Usage

```bash
# Monitor CPU usage
top -p $(docker-compose ps -q)

# Adjust CPU limits
services:
  ai-context-tracker:
    cpus: 2.0
```

## Monitoring and Maintenance

### Health Checks

```bash
# Automated health checks
curl -f http://localhost:8080/api/v1/health
curl -f http://localhost:5000/health
```

### Log Management

```bash
# View logs
docker-compose logs -f

# Log rotation (add to docker-compose.yml)
services:
  ai-context-tracker:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Database Maintenance

```bash
# Backup database
docker-compose exec postgres pg_dump -U tracker_user ai_context_tracker > backup.sql

# Restore database
docker-compose exec -T postgres psql -U tracker_user ai_context_tracker < backup.sql
```

### Updates and Upgrades

```bash
# Update to latest version
git pull origin main
make docker-build
docker-compose up -d

# Rollback if needed
git checkout previous-version
make docker-build
docker-compose up -d
```

## Security Considerations

### 1. Network Security

```yaml
# Use custom networks
networks:
  ai-tracker-network:
    driver: bridge
    internal: true
```

### 2. Environment Variables

```bash
# Use secrets management
docker secret create db_password db_password.txt
```

### 3. SSL/TLS Configuration

```bash
# Generate SSL certificates
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout server.key -out server.crt
```

## Scaling

### Horizontal Scaling

```yaml
# Scale services
docker-compose up -d --scale ai-context-tracker=3
```

### Load Balancing

```yaml
# Use nginx for load balancing
nginx:
  image: nginx:alpine
  ports:
    - "80:80"
  depends_on:
    - ai-context-tracker
```

## Support

For issues and questions:

1. Check the [troubleshooting section](#troubleshooting)
2. Review application logs
3. Create an issue in the GitHub repository
4. Contact the development team

## License

This project is licensed under the MIT License - see the LICENSE file for details.