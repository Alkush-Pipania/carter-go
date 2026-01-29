# Docker Compose Setup for Carter API

This guide explains how to run the Carter API project using Docker Compose.

## Prerequisites

- Docker installed
- Docker Compose installed
- `.env` file configured with your environment variables

## Configuration Files

- **`docker-compose.yml`** - Main configuration with all services (API, PostgreSQL, Redis, RabbitMQ)
- **`docker-compose.prod.yml`** - Production configuration with external services only
- **`docker-compose.override.yml`** - Override file for custom configurations
- **`.dockerignore`** - Files excluded from Docker build

## Quick Start

### Option 1: Development with All Services (Local)

Run everything including local PostgreSQL, Redis, and RabbitMQ:

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f carter-api

# Stop all services
docker-compose down
```

### Option 2: Production with External Services (EC2)

Run only the API with external PostgreSQL, Redis, and RabbitMQ:

```bash
# Build and start API only
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f carter-api

# Stop the service
docker-compose -f docker-compose.prod.yml down
```

## Environment Variables

Create a `.env` file in the project root. See [`.env.example`](.env.example) for reference.

### For Local Development (docker-compose.yml):

```bash
PORT=8080
ENV=development
LOG_LEVEL=info

# Use service names as hostnames
DB_URL=postgresql://postgres:postgres@postgres:5432/carter?sslmode=disable
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672
REDIS_ADDR=redis:6379

JWT_SECRET=your-secret-key-here
QUEUE_NAME=carter_queue
EXCHANGE_NAME=scarter.embedding
EXCHANGE_TYPE=direct
ROUTING_KEY=source.process
WORKER_COUNT=5

DO_REGION=blr1
DO_ENDPOINT=https://blr1.digitaloceanspaces.com
DO_ACCESS_KEY=your-access-key
DO_SECRET_KEY=your-secret-key
DO_BUCKET=your-bucket
PRESIGN_EXPIRY=15

REDIS_PASSWORD=
REDIS_DB=0
```

### For Production (docker-compose.prod.yml):

```bash
PORT=8080
ENV=production
LOG_LEVEL=info

# Use external service endpoints
DB_URL=postgres://avnadmin:password@your-db-host:port/dbname?sslmode=require
RABBITMQ_URL=amqp://guest:guest@localhost:5672
REDIS_ADDR=localhost:6379

JWT_SECRET=your-production-secret-key
QUEUE_NAME=carter_queue
EXCHANGE_NAME=scarter.embedding
EXCHANGE_TYPE=direct
ROUTING_KEY=source.process
WORKER_COUNT=5

DO_REGION=blr1
DO_ENDPOINT=https://blr1.digitaloceanspaces.com
DO_ACCESS_KEY=your-access-key
DO_SECRET_KEY=your-secret-key
DO_BUCKET=your-bucket
PRESIGN_EXPIRY=15

REDIS_PASSWORD=
REDIS_DB=0
```

## Common Commands

### Build and Start

```bash
# Development (all services)
docker-compose up -d --build

# Production (API only)
docker-compose -f docker-compose.prod.yml up -d --build
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f carter-api
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f rabbitmq
```

### Stop Services

```bash
# Stop and remove containers
docker-compose down

# Stop and remove containers with volumes
docker-compose down -v
```

### Restart Services

```bash
# Restart all services
docker-compose restart

# Restart specific service
docker-compose restart carter-api
```

### Execute Commands in Container

```bash
# Access API container shell
docker-compose exec carter-api sh

# Run migrations (if migrate tool is available)
docker-compose exec carter-api migrate -path migration -database "$DB_URL" up
```

## Health Check

Test if the API is running:

```bash
curl http://localhost:8080/health
```

Expected response: `OK`

## Services

### carter-api
- **Port**: 8080
- **Health Check**: `/health`
- **Environment**: Loaded from `.env` file

### postgres (Development only)
- **Port**: 5432
- **User**: postgres
- **Password**: postgres
- **Database**: carter
- **Data Volume**: `postgres-data`

### redis (Development only)
- **Port**: 6379
- **Data Volume**: `redis-data`

### rabbitmq (Development only)
- **Ports**: 5672 (AMQP), 15672 (Management UI)
- **User**: guest
- **Password**: guest
- **Management UI**: http://localhost:15672
- **Data Volume**: `rabbitmq-data`

## Troubleshooting

### Container Keeps Restarting

Check logs to identify the issue:

```bash
docker-compose logs carter-api
```

Common issues:
- Missing or incorrect environment variables
- Database connection failure
- Missing dependencies

### Database Connection Issues

For local development, ensure PostgreSQL is healthy:

```bash
docker-compose ps postgres
docker-compose logs postgres
```

For production, verify external database is accessible:

```bash
# Test connection from EC2
psql "postgresql://user:password@host:port/dbname?sslmode=require"
```

### RabbitMQ Connection Issues

Check RabbitMQ status:

```bash
docker-compose ps rabbitmq
docker-compose logs rabbitmq
```

Access RabbitMQ Management UI: http://localhost:15672

### Port Already in Use

If port 8080 is already in use, modify the port mapping in `docker-compose.yml`:

```yaml
ports:
  - "8081:8080"  # Use 8081 instead of 8080
```

## Production Deployment on EC2

1. **Copy project files to EC2**
2. **Create `.env` file** with production values
3. **Run with production config**:

```bash
docker-compose -f docker-compose.prod.yml up -d --build
```

4. **Verify deployment**:

```bash
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f carter-api
curl http://localhost:8080/health
```

5. **Configure EC2 Security Group** to allow inbound traffic on port 8080

## Data Persistence

Data is persisted in Docker volumes:
- `postgres-data` - PostgreSQL data
- `redis-data` - Redis data
- `rabbitmq-data` - RabbitMQ data

To remove all data:

```bash
docker-compose down -v
```

## Updating the Application

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build

# Or for production
docker-compose -f docker-compose.prod.yml up -d --build
```
