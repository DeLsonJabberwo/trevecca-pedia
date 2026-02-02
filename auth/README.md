# Trevecca-Pedia Authentication Service

MVP authentication service that issues JWTs and supports role-based access control.

## Features

- ✅ JWT-based authentication (HS256)
- ✅ Role-based access control (RBAC)
- ✅ Bcrypt password hashing
- ✅ PostgreSQL persistence
- ✅ Development seed user
- ✅ CORS support
- ✅ Health check endpoint

## Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 12+
- Docker (optional)

### Local Development

1. **Copy environment variables**

```bash
cp .env.example .env
```

Edit `.env` with your configuration.

2. **Set up database**

Create a PostgreSQL database and user:

```sql
CREATE DATABASE auth;
CREATE USER auth_user WITH PASSWORD 'auth_pass';
GRANT ALL PRIVILEGES ON DATABASE auth TO auth_user;
```

3. **Run migrations**

```bash
psql -U auth_user -d auth -f migrations/0001_init.sql
```

Or connect to the database and run the migration manually.

4. **Install dependencies**

```bash
go mod download
```

5. **Run the service**

```bash
# Source env vars
export $(cat .env | xargs)

# Run
go run cmd/auth/main.go
```

The service will start on port 8083 (or your configured PORT).

### Using Docker

1. **Build the image**

```bash
docker build -t trevecca-pedia-auth .
```

2. **Run the container**

```bash
docker run -p 8083:8083 \
  -e DATABASE_URL="postgres://auth_user:auth_pass@host.docker.internal:5432/auth?sslmode=disable" \
  -e JWT_SECRET="your-secret-key" \
  -e DEV_SEED=true \
  trevecca-pedia-auth
```

## API Endpoints

### Health Check

```
GET /healthz
```

Response:
```json
{
  "status": "ok"
}
```

### Login

```
POST /auth/login
Content-Type: application/json

{
  "email": "dev@trevecca.edu",
  "password": "devpass"
}
```

Success Response (200):
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "dev@trevecca.edu",
    "roles": ["reader", "contributor"]
  }
}
```

Error Response (401):
```json
{
  "error": "invalid credentials"
}
```

### Get Current User

```
GET /auth/me
Authorization: Bearer <token>
```

Success Response (200):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "dev@trevecca.edu",
  "roles": ["reader", "contributor"]
}
```

Error Response (401):
```json
{
  "error": "unauthorized"
}
```

## JWT Contract

The service issues JWT tokens with the following structure:

**Algorithm:** HS256

**Claims:**
- `sub` (string): User ID (UUID)
- `email` (string): User email
- `roles` (array): User roles
- `iss` (string): "trevecca-pedia-auth"
- `aud` (string): "trevecca-pedia"
- `iat` (number): Issued at timestamp
- `exp` (number): Expiration timestamp (default 24h from issue)

Example decoded JWT:
```json
{
  "sub": "123e4567-e89b-12d3-a456-426614174000",
  "email": "dev@trevecca.edu",
  "roles": ["reader", "contributor"],
  "iss": "trevecca-pedia-auth",
  "aud": "trevecca-pedia",
  "iat": 1640000000,
  "exp": 1640086400
}
```

## User Roles

| Role | Description |
|------|-------------|
| `reader` | Can browse and view wiki pages |
| `contributor` | Can create and edit wiki pages |
| `admin` | Elevated permissions (future use) |

## Configuration

All configuration is done via environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | `8083` | Server port |
| `DATABASE_URL` | Yes | - | PostgreSQL connection string |
| `JWT_SECRET` | Yes | - | Secret key for signing JWTs |
| `JWT_EXP_HOURS` | No | `24` | JWT expiration in hours |
| `CORS_ORIGINS` | No | `http://localhost:3000,http://localhost:5173` | Comma-separated allowed origins |
| `DEV_SEED` | No | `false` | Create dev user on startup |

## Development User

When `DEV_SEED=true`, the service creates a development user on startup:

- **Email:** `dev@trevecca.edu`
- **Password:** `devpass`
- **Roles:** `contributor`

⚠️ **Warning:** Only enable DEV_SEED in development environments!

## Database Schema

### Tables

**users**
- `id` (UUID, PK): User ID
- `email` (TEXT, UNIQUE): User email
- `password_hash` (TEXT): Bcrypt password hash
- `created_at` (TIMESTAMPTZ): Creation timestamp

**roles**
- `id` (SERIAL, PK): Role ID
- `name` (TEXT, UNIQUE): Role name

**user_roles**
- `user_id` (UUID, FK → users): User ID
- `role_id` (INT, FK → roles): Role ID
- Primary key: (user_id, role_id)

## Testing

### Manual Testing

```bash
# Health check
curl http://localhost:8083/healthz

# Login
curl -X POST http://localhost:8083/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@trevecca.edu","password":"devpass"}'

# Get current user (replace TOKEN)
curl http://localhost:8083/auth/me \
  -H "Authorization: Bearer TOKEN"
```

### Unit Tests

Run tests with:

```bash
go test ./...
```

## Integration with Other Services

### API Layer

The API layer should validate tokens by:

1. Extracting the token from `Authorization: Bearer <token>`
2. Sending a request to `GET /auth/me` with the token
3. Using the returned user info and roles for authorization

Alternatively, the API layer can validate tokens locally using the same JWT secret and validation logic.

### Wiki Service

The wiki service should enforce authorization based on roles:

- **Reader role**: Can view pages
- **Contributor role**: Can create/edit pages
- **Admin role**: Full access (future)

## Error Handling

All errors return JSON responses with appropriate HTTP status codes:

- `400 Bad Request`: Invalid request format
- `401 Unauthorized`: Invalid credentials or token
- `500 Internal Server Error`: Server error

Example error response:
```json
{
  "error": "error description"
}
```

## Security Notes

- Passwords are hashed using bcrypt (cost factor 12)
- JWT tokens expire after configured hours (default 24h)
- CORS is configured to only allow specified origins
- Database connections use connection pooling
- No sensitive data is logged

## Architecture

```
auth/
├── cmd/auth/          # Application entrypoint
├── internal/
│   ├── auth/          # JWT and password handling
│   ├── config/        # Configuration loading
│   ├── http/          # HTTP handlers and routing
│   └── store/         # Database operations
├── migrations/        # SQL migrations
├── Dockerfile         # Container build
├── go.mod            # Go dependencies
└── README.md         # This file
```

## Future Enhancements (Post-MVP)

- Microsoft SSO integration
- Password recovery flows
- User self-registration
- Refresh tokens
- Token revocation
- Rate limiting
- Audit logging
- More comprehensive tests

## License

Part of the Trevecca-Pedia project.
