# Integration Guide: Using Auth Service in Other Services

This guide shows how to integrate JWT authentication in the API Layer and Wiki Service.

## Overview

The auth service issues JWT tokens that other services can validate either by:
1. **Calling `/auth/me`** - Simpler, but adds network latency
2. **Validating tokens locally** - Faster, requires sharing JWT secret

For the MVP, we'll use **local validation** for better performance.

## Shared JWT Validation Package

Create a shared package that can be used by all services.

### Option 1: Copy the JWT validation code

You can copy the `auth/internal/auth/jwt.go` file to each service and validate tokens locally.

### Option 2: Call the auth service

Make an HTTP request to `GET /auth/me` with the token to validate it.

## Example: API Layer Integration

Here's how to add JWT middleware to the API Layer:

### 1. Add JWT validation to API Layer

Create `api-layer/internal/auth/jwt.go`:

```go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret   []byte
	issuer   string
	audience string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret:   []byte(secret),
		issuer:   "trevecca-pedia-auth",
		audience: "trevecca-pedia",
	}
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.Issuer != j.issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	if len(claims.Audience) == 0 || claims.Audience[0] != j.audience {
		return nil, fmt.Errorf("invalid audience")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
```

### 2. Add middleware to API Layer

Create `api-layer/middleware/auth.go`:

```go
package middleware

import (
	"api-layer/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Next()
	}
}

func RequireRole(jwtService *auth.JWTService, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure they're authenticated
		AuthRequired(jwtService)(c)
		
		if c.IsAborted() {
			return
		}

		// Check if they have the required role
		roles, exists := c.Get("user_roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		userRoles := roles.([]string)
		hasRole := false
		for _, role := range userRoles {
			if role == requiredRole || role == "admin" {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
```

### 3. Update API Layer main.go

```go
package main

import (
	"api-layer/handlers/wiki"
	"api-layer/internal/auth"
	"api-layer/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	// Initialize JWT service with same secret as auth service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-in-production-please" // dev default
	}
	jwtService := auth.NewJWTService(jwtSecret)

	// Public routes (no auth required)
	r.GET("/v1/wiki/pages", wiki.GetPages)
	r.GET("/v1/wiki/pages/:id", wiki.GetPage)
	r.GET("/v1/wiki/pages/:id/revisions", wiki.GetPageRevisions)
	r.GET("/v1/wiki/pages/:id/revisions/:rev", wiki.GetPageRevision)

	// Protected routes (contributor role required)
	protected := r.Group("/v1/wiki")
	protected.Use(middleware.RequireRole(jwtService, "contributor"))
	{
		protected.POST("/pages/new", wiki.PostNewPage)
		protected.POST("/pages/:id/revisions", wiki.PostPageRevision)
	}

	r.Run(":2745")
}
```

### 4. Add JWT_SECRET to environment

Add to your API Layer's `.env` or environment:

```bash
JWT_SECRET=dev-secret-key-change-in-production-please
```

**IMPORTANT:** This must be the SAME secret used by the auth service!

## Testing Protected Endpoints

### 1. Get a token

```bash
TOKEN=$(curl -s -X POST http://localhost:8083/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@trevecca.edu","password":"devpass"}' | jq -r .accessToken)

echo "Token: $TOKEN"
```

### 2. Call protected endpoint

```bash
# This should work (with valid token)
curl -X POST http://localhost:2745/v1/wiki/pages/new \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: multipart/form-data" \
  -F "slug=test-page" \
  -F "name=Test Page" \
  -F "author=dev@trevecca.edu" \
  -F "new_page=@test.md"

# This should fail (no token)
curl -X POST http://localhost:2745/v1/wiki/pages/new \
  -H "Content-Type: multipart/form-data" \
  -F "slug=test-page" \
  -F "name=Test Page"
```

## Role-Based Access Control

The MVP defines these roles:

| Role | Can Read | Can Create/Edit |
|------|----------|-----------------|
| reader | ✅ | ❌ |
| contributor | ✅ | ✅ |
| admin | ✅ | ✅ |

### Checking Roles in Handlers

```go
func (h *Handler) SomeProtectedAction(c *gin.Context) {
	// Get user info from context (set by middleware)
	userID, _ := c.Get("user_id")
	userEmail, _ := c.Get("user_email")
	userRoles, _ := c.Get("user_roles")

	// Use user info
	log.Printf("Action performed by user: %v (%v)", userEmail, userID)

	// Check specific role if needed
	roles := userRoles.([]string)
	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	// Your handler logic...
}
```

## Wiki Service Integration

The Wiki service should also validate tokens for write operations:

### Update Wiki handlers

```go
// In wiki/cmd/main.go

// Public routes - no auth
r.GET("/pages", handleGetPages)
r.GET("/pages/:id", handleGetPage)
r.GET("/pages/:id/revisions", handleGetRevisions)
r.GET("/pages/:id/revisions/:rev", handleGetRevision)

// Protected routes - require auth
authGroup := r.Group("/")
authGroup.Use(middleware.AuthRequired(jwtService))
{
	authGroup.POST("/pages/new", handlePostNewPage)
	authGroup.POST("/pages/:id/revisions", handlePostRevision)
}
```

## Frontend Integration

### 1. Login Flow

```javascript
// Login
async function login(email, password) {
  const response = await fetch('http://localhost:8083/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  
  const data = await response.json();
  
  if (response.ok) {
    // Store token in localStorage or sessionStorage
    localStorage.setItem('accessToken', data.accessToken);
    localStorage.setItem('user', JSON.stringify(data.user));
    return data;
  } else {
    throw new Error(data.error);
  }
}
```

### 2. Using Token in Requests

```javascript
// Get stored token
function getAuthToken() {
  return localStorage.getItem('accessToken');
}

// Make authenticated request
async function createWikiPage(pageData) {
  const token = getAuthToken();
  
  const response = await fetch('http://localhost:2745/v1/wiki/pages/new', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(pageData)
  });
  
  return response.json();
}

// Check if user is authenticated
function isAuthenticated() {
  const token = getAuthToken();
  if (!token) return false;
  
  // Optionally validate token by calling /auth/me
  return true;
}

// Logout
function logout() {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('user');
}
```

### 3. Automatic Token Refresh (Future)

For now, tokens expire after 24 hours. In the future, you can:
1. Implement refresh tokens
2. Show login modal when token expires
3. Automatically renew tokens before expiration

## Environment Setup Checklist

For all services to work together:

- [ ] Auth service running on port 8083
- [ ] Auth database running (auth-db on port 5433)
- [ ] Same JWT_SECRET in auth service and API layer
- [ ] API layer validates tokens using shared secret
- [ ] Wiki service validates tokens for write operations
- [ ] Frontend stores and sends tokens in Authorization header

## Troubleshooting

### "Invalid token" errors

1. **Check JWT_SECRET**: Must be identical in auth service and API layer
2. **Check token format**: Should be `Bearer <token>`, not just `<token>`
3. **Check expiration**: Tokens expire after JWT_EXP_HOURS (default 24h)
4. **Check issuer/audience**: Must match "trevecca-pedia-auth" and "trevecca-pedia"

### "Authorization required" errors

1. **Check Authorization header**: Must be present and correctly formatted
2. **Check CORS**: Ensure CORS_ORIGINS includes your frontend origin
3. **Check token in localStorage**: May have been cleared or expired

### Role permission errors

1. **Check user roles**: Call `/auth/me` to see what roles the user has
2. **Check middleware order**: AuthRequired must come before RequireRole
3. **Check role name spelling**: "contributor" not "contributer"

## Security Best Practices

1. **Never log tokens**: They're sensitive credentials
2. **Use HTTPS in production**: Tokens can be intercepted over HTTP
3. **Validate on every request**: Don't trust client-side validation
4. **Short token expiration**: 24h is reasonable for MVP, consider shorter for production
5. **Rotate JWT_SECRET**: Change periodically, especially if compromised
6. **Rate limit auth endpoints**: Prevent brute force attacks
7. **Use secure token storage**: httpOnly cookies are safer than localStorage

## Next Steps

1. Implement token refresh mechanism
2. Add password reset flow
3. Implement Microsoft SSO
4. Add audit logging for auth events
5. Set up monitoring and alerts
6. Implement rate limiting
