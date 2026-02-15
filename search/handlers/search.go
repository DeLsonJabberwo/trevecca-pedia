package handlers

import (
	"fmt"
	"net/http"
	"search/service"

	"github.com/gin-gonic/gin"
)

var searchService *service.SearchService

func SetSearchService(s *service.SearchService) {
	searchService = s
}

type SearchResponse struct {
	Total   int      `json:"total"`
	Results []string `json:"results"`
}

func SearchHandler(c *gin.Context) {
	query := c.Query("q")

	searchResults, err := searchService.Search(query, 0, 10)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("err: %s", err))
		return
	}

	results := make([]string, 0, len(searchResults.Hits))
	for _, hit := range searchResults.Hits {
		results = append(results, hit.ID)
	}

	response := SearchResponse{
		Total:   int(searchResults.Total),
		Results: results,
	}

	c.JSON(http.StatusOK, response)
}

func ReindexHandler(c *gin.Context) {
	err := searchService.IndexAll("")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("err: %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reindex completed successfully"})
}
