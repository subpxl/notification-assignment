# Automated Messaging Service

A Go-based messaging service that automatically sends messages from a PostgreSQL database to a webhook endpoint, with Redis caching and a built-in scheduler.


## Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)

## Quick Start

### Using Docker Compose (Recommended)


2. Start all services:
```bash
docker-compose up --build
```

The application will be available at:
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

### Using Makefile

```bash

# Start services in detached mode
make up

# Stop services
make down

# Build the application locally
make build

# Check Redis cached messages
make redis-keys

# Check psql messages
make psql
```

## Configuration

Configure the application via environment variables in [docker-compose.yml](docker-compose.yml):

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `postgres` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `insider_messaging` |
| `REDIS_HOST` | Redis host | `redis` |
| `REDIS_PORT` | Redis port | `6379` |
| `WEBHOOK_URL` | Target webhook URL | - |
| `AUTH_KEY` | Authentication key for webhook | - |
| `INTERVAL_SEC` | Scheduler interval in seconds | `120` |
| `BATCH_SIZE` | Messages per batch | `2` |

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

### Control Scheduler
```bash
# Start message sending
curl -X POST http://localhost:8080/api/v1/control/start

# Stop message sending
curl -X POST http://localhost:8080/api/v1/control/stop
```

### Retrieve Sent Messages
```bash
curl http://localhost:8080/api/v1/sent-messages
```

## Webhook Request Format

The service sends messages to the configured webhook endpoint with:

**Headers:**
- `Content-Type: application/json`
- `x-ins-auth-key: <AUTH_KEY>`

**Payload:**
```json
{
  "to": "+905551111111",
  "content": "Your message content"
}
```

## Redis Cache

Messages are cached in Redis after successful sending. View cached data:

```bash
# Access Redis CLI
docker exec -it insider-assignment-redis-1 redis-cli

# List all keys
KEYS *

# View specific message
GET sent_message:1

# View all sent messages
KEYS sent_message:*
```

## Database

PostgreSQL is initialized with the schema from [init-sql/init-db.sql](init-sql/init-db.sql). The scheduler automatically processes unsent messages.

## Development

### Regenerate Swagger Documentation

```bash
make swagger
# or
swag init -g cmd/main.go
```

