package wiki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"web/config"
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

	component := wikipages.Page(page)
    component.Render(context.Background(), c.Writer)

}
