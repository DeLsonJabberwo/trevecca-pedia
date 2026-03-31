package mod

import (
	"api-layer/config"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func GetStatuses(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	statusURL, err := url.Parse(fmt.Sprintf("%s/statuses", config.ModServiceURL))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}
	q := statusURL.Query()
	q.Set("username", username)
	statusURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, statusURL.String(), nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Forward Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user statuses."})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}

func GetFlaggedUsers(c *gin.Context) {
	ind := c.DefaultQuery("ind", "")
	count := c.DefaultQuery("count", "")
	url, err := url.Parse(fmt.Sprintf("%s/flagged-users", config.ModServiceURL))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
	}
	q := url.Query()
	q.Set("ind", ind)
	q.Set("count", count)
	url.RawQuery = q.Encode()
	res, err := http.Get(url.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user statuses."})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}

func GetSilencedUsers(c *gin.Context) {
	ind := c.DefaultQuery("ind", "")
	count := c.DefaultQuery("count", "")
	url, err := url.Parse(fmt.Sprintf("%s/silenced-users", config.ModServiceURL))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
	}
	q := url.Query()
	q.Set("ind", ind)
	q.Set("count", count)
	url.RawQuery = q.Encode()
	res, err := http.Get(url.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user statuses."})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}
