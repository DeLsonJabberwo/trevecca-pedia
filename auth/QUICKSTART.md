# Auth Service Quick Start

Get the authentication service up and running in 5 minutes.

## Prerequisites

- Docker & Docker Compose
- Go 1.22+
- curl (for testing)

## Setup

### 1. Start the Database

From the `wiki-db` directory:

```bash
cd wiki-db
docker-compose up -d auth-db
```

Wait for the database to be healthy (~10 seconds):

```bash
docker-compose ps
```

### 2. Run the Auth Service

From the `auth` directory:

```bash
cd ../auth

# Install dependencies
go mod download

# Run the service
go run cmd/auth/main.go
```

The service will:
- Connect to the database
- Create a dev user: `dev@trevecca.edu` / `devpass`
- Start on port 8083

### 3. Test It

```bash
# Health check
curl http://localhost:8083/healthz

# Login
curl -X POST http://localhost:8083/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@trevecca.edu","password":"devpass"}'
```

You should receive a JWT token and user info.

## Alternative: Docker Only

Run everything with Docker:

```bash
cd auth
docker-compose up
```

This starts both the database and auth service.

## Full Test Suite

Run the automated test suite:

```bash
cd auth
./test-auth.sh
```

## Documentation

- [README.md](./README.md) - Full documentation
- [SETUP.md](./SETUP.md) - Detailed setup guide
- [INTEGRATION.md](./INTEGRATION.md) - Integration with other services

## Troubleshooting

**Database connection fails:**
```bash
# Check if auth-db is running
cd ../wiki-db
docker-compose ps

# Check logs
docker logs auth-db
```

**Port 8083 already in use:**
```bash
# Change PORT in auth/.env
PORT=8084
```

**Dev user not created:**
```bash
# Check DEV_SEED in auth/.env
DEV_SEED=true
```

## Next Steps

1. Read [INTEGRATION.md](./INTEGRATION.md) to integrate with API Layer
2. Update API Layer to validate JWT tokens
3. Add authentication to Wiki Service
4. Build login UI in Web Service
