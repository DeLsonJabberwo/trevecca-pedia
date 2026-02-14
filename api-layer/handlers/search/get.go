package search

import (
	"api-layer/config"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchRequest(c *gin.Context) {
	query := c.Query("q")
	url := fmt.Sprintf("%s/search?q=%s", config.SearchServiceURL, query)
	resp, err := http.Get(url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	c.Data(http.StatusOK, resp.Header.Get("Content-Type"), body)
}
