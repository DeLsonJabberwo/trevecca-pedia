package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"web/config"

	"github.com/gin-gonic/gin"
)

func GetSearchPage(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		// Render empty search page
		return
	}
	searchResp, err := http.Get(fmt.Sprintf("%s/search?q=%s", config.SearchURL, url.QueryEscape(query)))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't fetch search results",
		})
	}
	defer searchResp.Body.Close()

	searchBody, err := io.ReadAll(searchResp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't fetch search results",
		})
	}
	var slugs []string
	err = json.Unmarshal(searchBody, &slugs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't fetch search results",
		})
	}

	slugsParam := strings.Join(slugs, ",")
	wikiResp, err := http.Get(fmt.Sprintf("%s/pages?slugs=%s", config.WikiURL, url.QueryEscape(slugsParam)))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't fetch search results",
		})
	}
	defer wikiResp.Body.Close()

	wikiBody, err := io.ReadAll(wikiResp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't fetch search results",
		})
	}
	// this part is in another branch right now. oops


}

