package auth

import (
	"bytes"
	"io"
	"net/http"
	"web/config"

	"github.com/gin-gonic/gin"
)

// PostLogin proxies POST /auth/login → API layer POST /v1/auth/login
func PostLogin(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to read request"})
		return
	}

	res, err := http.Post(config.AuthURL+"/login", "application/json", bytes.NewReader(body))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "auth service unavailable"})
		return
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), respBody)
}

// PostRegister proxies POST /auth/register → API layer POST /v1/auth/register
func PostRegister(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to read request"})
		return
	}

	res, err := http.Post(config.AuthURL+"/register", "application/json", bytes.NewReader(body))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "auth service unavailable"})
		return
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), respBody)
}

// GetMe proxies GET /auth/me → API layer GET /v1/auth/me
// The Authorization header is forwarded so the auth service can validate the JWT.
func GetMe(c *gin.Context) {
	req, err := http.NewRequest(http.MethodGet, config.AuthURL+"/me", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "auth service unavailable"})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}
