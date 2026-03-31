package mod

import (
	"api-layer/config"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostFlagUser(c *gin.Context) {
	// Read the form data from the request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	defer c.Request.Body.Close()

	// Forward to moderation service
	modURL := fmt.Sprintf("%s/flag-user", config.ModServiceURL)
	req, err := http.NewRequest(http.MethodPost, modURL, bytes.NewReader(body))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Forward Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to flag user"})
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), resBody)
}

func PostUnflagUser(c *gin.Context) {
	// Read the form data from the request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	defer c.Request.Body.Close()

	// Forward to moderation service
	modURL := fmt.Sprintf("%s/unflag-user", config.ModServiceURL)
	req, err := http.NewRequest(http.MethodPost, modURL, bytes.NewReader(body))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Forward Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to unflag user"})
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), resBody)
}

func PostSilenceUser(c *gin.Context) {
	// Read the form data from the request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	defer c.Request.Body.Close()

	// Forward to moderation service
	modURL := fmt.Sprintf("%s/silence-user", config.ModServiceURL)
	req, err := http.NewRequest(http.MethodPost, modURL, bytes.NewReader(body))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Forward Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to silence user"})
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), resBody)
}

func PostUnsilenceUser(c *gin.Context) {
	// Read the form data from the request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	defer c.Request.Body.Close()

	// Forward to moderation service
	modURL := fmt.Sprintf("%s/unsilence-user", config.ModServiceURL)
	req, err := http.NewRequest(http.MethodPost, modURL, bytes.NewReader(body))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Forward Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsilence user"})
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), resBody)
}
