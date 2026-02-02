# Authentication Service - Setup Guide

This guide walks you through setting up and running the Trevecca-Pedia authentication service.

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 12+ (or Docker to run the database)
- Docker (optional, for containerized database and service)

## Quick Start (Using Docker for Database)

### 1. Start the Database

From the **wiki-db** directory (one level up from auth):

```bash
cd ../wiki-db
docker-compose up -d auth-db
```

This will:
- Start a PostgreSQL container on port 5433
- Create the `auth` database
- Create the `auth_user` user
- Run the migrations automatically (from `init-auth/0001_init.sql`)

Check if the database is ready:

```bash
docker-compose ps
```

You should see `auth-db` with status `healthy`.

### 2. Configure the Auth Service

The auth service comes with a pre-configured `.env` file for local development. Review and modify if needed:

```bash
cd ../auth
cat .env
```

Key settings:
- `PORT=8083` - Service port
- `DATABASE_URL` - Connection to auth-db on port 5433
- `JWT_SECRET` - Secret for signing tokens (change in production!)
- `DEV_SEED=true` - Creates a test user on startup

### 3. Install Dependencies

```bash
make deps
# or
go mod download
```

### 4. Run the Service

Option A: Direct run
```bash
make run
# or
go run cmd/auth/main.go
```

Option B: With hot reload (requires [air](https://github.com/air-verse/air))
```bash
# Install air if you don't have it
go install github.com/air-verse/air@latest

# Run with hot reload
make dev
# or
air
```

The service will start on port 8083.

You should see output like:
```
Database connection established
⚠️  DEV_SEED is enabled - creating development user
✓ Development user created/verified: dev@trevecca.edu / devpass
Starting auth service on port 8083
```

### 5. Test the Service

Health check:
```bash
curl http://localhost:8083/healthz
```

Login with dev user:
```bash
curl -X POST http://localhost:8083/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@trevecca.edu","password":"devpass"}'
```

You should receive a response with an access token and user info.

Get current user (replace TOKEN with the token from login):
```bash
curl http://localhost:8083/auth/me \
  -H "Authorization: Bearer TOKEN"
```

## Alternative Setup (Manual Database)

If you prefer to set up PostgreSQL manually without Docker:

### 1. Create Database and User

```sql
-- Connect to PostgreSQL as superuser
psql -U postgres

-- Create user
CREATE USER auth_user WITH PASSWORD 'authpass';

-- Create database
CREATE DATABASE auth OWNER auth_user;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE auth TO auth_user;

-- Connect to the new database
\c auth

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO auth_user;
```

### 2. Run Migrations

```bash
psql -U auth_user -d auth -f migrations/0001_init.sql
```

### 3. Update .env

Update the `DATABASE_URL` in `.env` to match your local PostgreSQL setup:

```
DATABASE_URL=postgres://auth_user:authpass@localhost:5432/auth?sslmode=disable
```

Note: Port 5432 if you're running PostgreSQL directly (not in Docker).

### 4. Run the Service

Follow steps 3-5 from the Docker setup above.

## Running with Docker

### Build the Docker Image

```bash
make docker-build
# or
docker build -t trevecca-pedia-auth:latest .
```

### Run the Container

```bash
docker run -p 8083:8083 \
  -e DATABASE_URL="postgres://auth_user:authpass@host.docker.internal:5433/auth?sslmode=disable" \
  -e JWT_SECRET="dev-secret-key-change-in-production-please" \
  -e DEV_SEED=true \
  trevecca-pedia-auth:latest
```

Note: Use `host.docker.internal` to connect to databases on your host machine from Docker.

## Development Workflow

### Hot Reload with Air

For development, use Air for automatic reloading on code changes:

```bash
make dev
```

Air will watch for Go file changes and automatically rebuild and restart the service.

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Formatting

```bash
make fmt
```

### Building Binary

```bash
# Build to bin/auth
make build

# Run the binary
./bin/auth
```

## Verifying the Setup

### 1. Check Database Connection

```bash
# Connect to the database
docker exec -it auth-db psql -U auth_user -d auth

# List tables
\dt

# You should see: users, roles, user_roles

# Check roles
SELECT * FROM roles;

# Exit
\q
```

### 2. Check Dev User

```sql
-- Connect to database
docker exec -it auth-db psql -U auth_user -d auth

-- Check if dev user exists
SELECT u.id, u.email, u.created_at, r.name as role
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.email = 'dev@trevecca.edu';
```

### 3. Manual API Testing

Create a file `test-auth.sh`:

```bash
#!/bin/bash

# Test health check
echo "Testing health check..."
curl -s http://localhost:8083/healthz | jq .
echo ""

# Test login
echo "Testing login..."
TOKEN=$(curl -s -X POST http://localhost:8083/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@trevecca.edu","password":"devpass"}' | jq -r .accessToken)

echo "Received token: ${TOKEN:0:50}..."
echo ""

# Test me endpoint
echo "Testing /auth/me..."
curl -s http://localhost:8083/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Make it executable and run:
```bash
chmod +x test-auth.sh
./test-auth.sh
```

## Troubleshooting

### Database Connection Fails

**Error:** `Failed to connect to database`

Solutions:
1. Check if the database container is running: `docker-compose ps`
2. Check if the port is correct (5433 for Docker, 5432 for local)
3. Verify credentials in `.env` match `wiki-db/.env`
4. Check if another process is using port 5433: `lsof -i :5433`

### Dev User Not Created

**Issue:** DEV_SEED is true but user not created

Solutions:
1. Check the logs for error messages
2. Ensure migrations ran: `docker exec -it auth-db psql -U auth_user -d auth -c "\dt"`
3. Try creating user manually via SQL
4. Check if roles table has data: `SELECT * FROM roles;`

### Port Already in Use

**Error:** `bind: address already in use`

Solutions:
1. Check what's using port 8083: `lsof -i :8083`
2. Kill the process or change PORT in `.env`
3. Update `.air.toml` if using Air with a different port

### Invalid Token Errors

**Issue:** Token validation fails

Solutions:
1. Ensure JWT_SECRET is the same when generating and validating
2. Check token hasn't expired (default 24h)
3. Verify token format: should be `Bearer <token>`
4. Check server logs for specific validation errors

### Migration Issues

**Issue:** Tables not created

Solutions:
1. Manually run migrations:
   ```bash
   docker exec -i auth-db psql -U auth_user -d auth < migrations/0001_init.sql
   ```
2. Check Docker logs: `docker logs auth-db`
3. Drop and recreate database if needed (dev only!)

## Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| PORT | No | 8083 | HTTP server port |
| DATABASE_URL | Yes | - | PostgreSQL connection string |
| JWT_SECRET | Yes | - | Secret for signing JWTs (min 32 chars recommended) |
| JWT_EXP_HOURS | No | 24 | Token expiration in hours |
| CORS_ORIGINS | No | localhost:3000,5173,8080 | Comma-separated allowed CORS origins |
| DEV_SEED | No | false | Create dev user on startup (dev only!) |

## Next Steps

Once the auth service is running:

1. **Integrate with API Layer**: Update the API layer to validate JWT tokens
2. **Update Wiki Service**: Add authentication checks for create/edit operations
3. **Frontend Integration**: Add login UI and token management
4. **Test Full Flow**: Create/edit wiki pages with authentication

## Useful Commands

```bash
# View auth service logs
docker logs -f auth-db

# Connect to auth database
docker exec -it auth-db psql -U auth_user -d auth

# Restart database
cd ../wiki-db && docker-compose restart auth-db

# Stop database
docker-compose stop auth-db

# Remove database (WARNING: deletes all data)
docker-compose down -v auth-db

# Build and run everything
make clean && make build && make run

# Run tests continuously
watch -n 1 make test
```

## Production Considerations

Before deploying to production:

1. ✅ Change JWT_SECRET to a strong random value
2. ✅ Set DEV_SEED=false
3. ✅ Use proper DATABASE_URL with SSL (sslmode=require)
4. ✅ Set up proper CORS_ORIGINS (no wildcards)
5. ✅ Use environment-based config (not .env files)
6. ✅ Set up monitoring and logging
7. ✅ Use a reverse proxy (nginx/traefik)
8. ✅ Implement rate limiting
9. ✅ Regular database backups
10. ✅ Security audit of JWT configuration

## Support

For issues or questions:
1. Check the logs: `docker logs auth-db` or service console output
2. Review the main README.md
3. Check the arc42 documentation
4. Ask the team in your project channel
