package wiki

import (
	"api-layer/config"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPages(c *gin.Context) {
	catQuery := c.DefaultQuery("category", "")
	ind, err := strconv.Atoi(c.DefaultQuery("index", "0"))
	if err != nil {
		ind = 0
	}
	num, err := strconv.Atoi(c.DefaultQuery("num", "10"))
	if err != nil {
		num = 10
	}

	res, err := http.Get(fmt.Sprintf("%s/pages?category=%s&index=%d&num=%d",
							config.WikiServiceURL, catQuery, ind, num))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages."})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
        return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}

func GetPage(c *gin.Context) {
	id := c.Param("id")
	res, err := http.Get(fmt.Sprintf("%s/pages/%s", config.WikiServiceURL, id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages."})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
        return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}

func GetPageRevisions(c *gin.Context) {
	id := c.Param("id")
	ind, err := strconv.Atoi(c.DefaultQuery("index", "0"))
	if err != nil {
		ind = 0
	}
	num, err := strconv.Atoi(c.DefaultQuery("num", "10"))
	if err != nil {
		num = 10
	}

	res, err := http.Get(fmt.Sprintf("%s/pages/%s/revisions?index=%d&num=%d",
							config.WikiServiceURL, id, ind, num))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages."})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
        return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}

func GetPageRevision(c *gin.Context) {
	id := c.Param("id")
	revId := c.Param("rev")
	res, err := http.Get(fmt.Sprintf("%s/pages/%s/revisions/%s", config.WikiServiceURL, id, revId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages."})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
        return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}
