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
		fmt.Printf("%s\n", body)
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Couldn't parse json from API layer: %w\n", err))
	}

	entryContent := wikipages.WikiEntryContent(page)
	component := components.Page(page.Name, entryContent)
	component.Render(context.Background(), c.Writer)

}

func GetHome(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	pages, err := getPages()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Couldn't fetch page info: %w\n", err))
	}
	homeComp := components.HomeContent(pages)
	page := components.Page("TreveccaPedia", homeComp)
	page.Render(context.Background(), c.Writer)
}

func getPages() ([]utils.Page, error) {
	resp, err := http.Get(fmt.Sprintf("%s/pages", config.WikiURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pages []utils.Page
	err = json.Unmarshal(body, &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

