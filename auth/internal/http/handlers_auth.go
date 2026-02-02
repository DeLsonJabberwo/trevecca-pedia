package http

import (
	"log"
	"net/http"

	"auth/internal/auth"
	"auth/internal/store"

	"github.com/gin-gonic/gin"
)

// AuthHandlers handles authentication endpoints
type AuthHandlers struct {
	store      *store.Store
	jwtService *auth.JWTService
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(store *store.Store, jwtService *auth.JWTService) *AuthHandlers {
	return &AuthHandlers{
		store:      store,
		jwtService: jwtService,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	AccessToken string                `json:"accessToken"`
	User        *store.UserWithRoles `json:"user"`
}

// Login handles user login
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	// Get user by email
	user, err := h.store.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		log.Printf("Login failed for %s: user not found", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify password
	if err := auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		log.Printf("Login failed for %s: invalid password", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Get user roles
	roles, err := h.store.GetUserRoles(c.Request.Context(), user.ID)
	if err != nil {
		log.Printf("Error getting roles for user %s: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user.ID, user.Email, roles)
	if err != nil {
		log.Printf("Error generating token for user %s: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Printf("User %s logged in successfully", user.Email)

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: token,
		User: &store.UserWithRoles{
			ID:    user.ID,
			Email: user.Email,
			Roles: roles,
		},
	})
}

// Me handles getting current user info
func (h *AuthHandlers) Me(c *gin.Context) {
	// Get claims from context (set by AuthMiddleware)
	claimsValue, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsValue.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return user info from claims
	c.JSON(http.StatusOK, store.UserWithRoles{
		ID:    claims.UserID,
		Email: claims.Email,
		Roles: claims.Roles,
	})
}

// HealthCheck handles health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
