package wiki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"web/config"
	"web/templates/components"
	wikipages "web/templates/wiki-pages"
	"web/utils"

	"github.com/gin-gonic/gin"
)

func PostEditPage(c *gin.Context) {
	id := c.Param("id")

	// Step 1 — fetch the page to get its UUID
	resp, err := http.Get(fmt.Sprintf("%s/pages/%s", config.WikiURL, id))
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var page utils.Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Step 2 — read the textarea content
	content := c.PostForm("content")
	if content == "" {
		editContent := wikipages.WikiEditContent(page, "Content cannot be empty.")
		component := components.Page("Editing: "+page.Name, editContent)
		component.Render(context.Background(), c.Writer)
		return
	}

	// Step 3 — forward to api-layer
	editURL := fmt.Sprintf("%s/pages/%s/edit", config.WikiURL, page.UUID)

	formData := url.Values{}
	formData.Set("content", content)
	formData.Set("author", c.PostForm("author"))

	editResp, err := http.Post(
		editURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		// Network error — wiki service unreachable
		editContent := wikipages.WikiEditContent(page, "Unable to save changes. The wiki service is unreachable.")
		component := components.Page("Editing: "+page.Name, editContent)
		component.Render(context.Background(), c.Writer)
		return
	}
	defer editResp.Body.Close()

	if editResp.StatusCode != http.StatusOK {
		// Wiki service returned an error — read the message
		respBody, _ := io.ReadAll(editResp.Body)
		errMsg := fmt.Sprintf("Unable to save changes. (status %d: %s)", editResp.StatusCode, string(respBody))
		editContent := wikipages.WikiEditContent(page, errMsg)
		component := components.Page("Editing: "+page.Name, editContent)
		component.Render(context.Background(), c.Writer)
		return
	}

	// Step 4 — success, redirect back to the page
	c.Redirect(http.StatusFound, fmt.Sprintf("/pages/%s?saved=true", id))
}