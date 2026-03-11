package wiki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"web/config"
	"web/templates/components"
	wikipages "web/templates/wiki-pages"
	"web/utils"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) {
	id := c.Param("id")
	resp, err := http.Get(fmt.Sprintf("%s/pages/%s", config.WikiURL, id))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Couldn't read http request: %w\n", err))
	}

	var page utils.Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Couldn't parse json from API layer: %w\n", err))
	}

	saved := c.Query("saved") == "true"
	entryContent := wikipages.WikiEntryContent(page, saved)
	component := components.Page(page.Name, entryContent)
	component.Render(context.Background(), c.Writer)

}

func GetHome(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	categories, err := getCategories()
	if err != nil {
		categories = []utils.Category{}
	}
	homeComp := components.HomeContent(categories)
	page := components.Page("TreveccaPedia", homeComp)
	page.Render(context.Background(), c.Writer)
}

func getPages() ([]utils.PageInfoPrev, error) {
	resp, err := http.Get(fmt.Sprintf("%s/pages", config.WikiURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pages []utils.PageInfoPrev
	err = json.Unmarshal(body, &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func getCategories() ([]utils.Category, error) {
	resp, err := http.Get(fmt.Sprintf("%s/categories?root=true", config.WikiURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var categories []utils.Category
	err = json.Unmarshal(body, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func GetEditPage(c *gin.Context) {
	id := c.Param("id")
	resp, err := http.Get(fmt.Sprintf("%s/pages/%s", config.WikiURL, id))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("couldn't read http response: %w", err))
		return
	}

	var page utils.Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("couldn't parse json from API layer: %w", err))
		return
	}

	editContent := wikipages.WikiEditContent(page, "")
	component := components.Page("Editing: "+page.Name, editContent)
	component.Render(context.Background(), c.Writer)
}

// PostPreview handles markdown preview requests from the editor
type PreviewRequest struct {
	Content string `json:"content"`
}

type PreviewResponse struct {
	HTML string `json:"html"`
}

func PostPreview(c *gin.Context) {
	var req PreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	html, err := utils.ToHTML(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render markdown"})
		return
	}

	c.JSON(http.StatusOK, PreviewResponse{HTML: html})
}
