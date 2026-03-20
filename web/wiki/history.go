package wiki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"web/config"
	"web/templates/components"
	wikipages "web/templates/wiki-pages"
	"web/utils"

	"github.com/gin-gonic/gin"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// GetPageHistory renders the split-view revision history page
func GetPageHistory(c *gin.Context) {
	id := c.Param("id")
	revId := c.Param("revId")

	// Get current page data
	page, err := fetchPageData(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Get revisions - fetch more to ensure we can find previous revision for older entries
	revisions, err := fetchRevisions(id, 0, 100)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Determine which revision to show
	var currentRevision utils.Revision
	var revisionNumber int

	if revId != "" {
		// Fetch specific revision
		currentRevision, err = fetchRevision(id, revId)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		// Find the revision number (position in list, with oldest as #1)
		for i, rev := range revisions {
			if rev.UUID == currentRevision.UUID {
				revisionNumber = len(revisions) - i
				break
			}
		}
	} else {
		// Show the most recent revision (first in the list since sorted newest first)
		if len(revisions) > 0 {
			// Fetch the full revision content since the list endpoint may not include it
			currentRevision, err = fetchRevision(id, revisions[0].UUID.String())
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			revisionNumber = len(revisions)
		}
	}

	// Get previous revision for diff highlighting
	var previousRevision *utils.Revision
	for i, rev := range revisions {
		if rev.UUID == currentRevision.UUID && i < len(revisions)-1 {
			// Fetch the full previous revision content since the list doesn't include it
			prevRevId := revisions[i+1].UUID.String()
			prevRev, err := fetchRevision(id, prevRevId)
			if err == nil {
				previousRevision = &prevRev
			}
			break
		}
	}

	// Highlight changes and convert to HTML
	highlightedContent, hasChanges := highlightChanges(currentRevision.Content, previousRevision)

	// Check if HTMX request (for partial content update)
	if c.GetHeader("HX-Request") == "true" {
		// Return article content AND updated timeline selection
		// Article replaces #article-content via hx-target
		articleContent := wikipages.WikiHistoryArticle(page, currentRevision, highlightedContent, revisionNumber, hasChanges)
		articleContent.Render(context.Background(), c.Writer)

		// Timeline updates selection via hx-swap-oob
		timelineContent := wikipages.WikiHistoryTimeline(revisions, currentRevision.UUID.String(), len(revisions))
		timelineContent.Render(context.Background(), c.Writer)
		return
	}

	// Full page render
	historyContent := wikipages.WikiHistoryContent(page, revisions, currentRevision, highlightedContent, revisionNumber, hasChanges)
	component := components.Page(page.Name+" - Revision History", historyContent)
	component.Render(context.Background(), c.Writer)
}

// GetRevisionContent returns just the article content for HTMX swaps
func GetRevisionContent(c *gin.Context) {
	id := c.Param("id")
	revId := c.Param("revId")

	page, err := fetchPageData(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	revision, err := fetchRevision(id, revId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Get all revisions to find the number and previous revision
	revisions, err := fetchRevisions(id, 0, 100)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var revisionNumber int
	var previousRevision *utils.Revision
	for i, rev := range revisions {
		if rev.UUID == revision.UUID {
			revisionNumber = len(revisions) - i
			if i < len(revisions)-1 {
				// Fetch the full previous revision content since the list doesn't include it
				prevRevId := revisions[i+1].UUID.String()
				prevRev, err := fetchRevision(id, prevRevId)
				if err == nil {
					previousRevision = &prevRev
				}
			}
			break
		}
	}

	highlightedContent, hasChanges := highlightChanges(revision.Content, previousRevision)

	// Render both article and updated timeline selection
	// Article replaces #article-content, timeline updates selection via hx-swap-oob
	articleContent := wikipages.WikiHistoryArticle(page, revision, highlightedContent, revisionNumber, hasChanges)
	timelineContent := wikipages.WikiHistoryTimeline(revisions, revId, len(revisions))

	// First render article, then timeline with oob swap
	articleContent.Render(context.Background(), c.Writer)
	// Timeline will have hx-swap-oob="true" to update out-of-band
	timelineContent.Render(context.Background(), c.Writer)
}

// GetTimelinePartial returns more timeline items for infinite scroll
func GetTimelinePartial(c *gin.Context) {
	id := c.Param("id")

	indexStr := c.Query("index")
	index := 0
	if indexStr != "" {
		index, _ = strconv.Atoi(indexStr)
	}

	revisions, err := fetchRevisions(id, index, 20)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Get total count for numbering
	totalRevisions, _ := fetchRevisions(id, 0, 1000)
	totalCount := len(totalRevisions)

	timelineItems := wikipages.WikiHistoryTimelineItems(revisions, totalCount, index)
	timelineItems.Render(context.Background(), c.Writer)
}

// fetchPageData gets page data from API
func fetchPageData(id string) (utils.Page, error) {
	resp, err := http.Get(fmt.Sprintf("%s/pages/%s", config.WikiURL, id))
	if err != nil {
		return utils.Page{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.Page{}, err
	}

	var page utils.Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return utils.Page{}, err
	}

	return page, nil
}

// fetchRevisions gets revision list from API
func fetchRevisions(id string, index, count int) ([]utils.Revision, error) {
	url := fmt.Sprintf("%s/pages/%s/revisions?index=%d&count=%d", config.WikiURL, id, index, count)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var revisions []utils.Revision
	err = json.Unmarshal(body, &revisions)
	if err != nil {
		return nil, err
	}

	return revisions, nil
}

// fetchRevision gets a specific revision from API
func fetchRevision(id, revId string) (utils.Revision, error) {
	url := fmt.Sprintf("%s/pages/%s/revisions/%s", config.WikiURL, id, revId)
	resp, err := http.Get(url)
	if err != nil {
		return utils.Revision{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.Revision{}, err
	}

	var revision utils.Revision
	err = json.Unmarshal(body, &revision)
	if err != nil {
		return utils.Revision{}, err
	}

	return revision, nil
}

// highlightChanges compares content and highlights changes using word-level diff
// Returns the highlighted HTML content and a boolean indicating if there are changes
func highlightChanges(currentContent string, previousRevision *utils.Revision) (string, bool) {
	// Convert current content to HTML
	currentHTML, err := utils.ToHTML(currentContent)
	if err != nil {
		currentHTML = currentContent
	}

	if previousRevision == nil {
		// No previous revision, return HTML as-is with no changes
		return currentHTML, false
	}

	// Convert previous content to HTML
	previousHTML, err := utils.ToHTML(previousRevision.Content)
	if err != nil {
		previousHTML = previousRevision.Content
	}

	// If HTML outputs are identical, no changes to highlight
	if currentHTML == previousHTML {
		return currentHTML, false
	}

	// Use diffmatchpatch to find differences on HTML
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(previousHTML, currentHTML, true)

	// Clean up the diff to be more semantic (word-level rather than char-level)
	diffs = dmp.DiffCleanupSemantic(diffs)

	// Check if there are any actual changes (insertions or deletions)
	hasChanges := false
	for _, diff := range diffs {
		if diff.Type == diffmatchpatch.DiffInsert || diff.Type == diffmatchpatch.DiffDelete {
			hasChanges = true
			break
		}
	}

	if !hasChanges {
		return currentHTML, false
	}

	// Build result with highlighted changes
	var result strings.Builder

	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			// New content - highlight as addition
			highlighted := highlightInsert(diff.Text)
			result.WriteString(highlighted)
		case diffmatchpatch.DiffDelete:
			// Deleted content - don't show in current version
			continue
		case diffmatchpatch.DiffEqual:
			// Unchanged content
			result.WriteString(diff.Text)
		}
	}

	return result.String(), hasChanges
}

// blockLevelTags are HTML tags that create block-level elements
var blockLevelTags = []string{"<p", "</p>", "<div", "</div>", "<h1", "</h1>", "<h2", "</h2>",
	"<h3", "</h3>", "<h4", "</h4>", "<h5", "</h5>", "<h6", "</h6>",
	"<ul", "</ul>", "<ol", "</ol>", "<li", "</li>", "<blockquote", "</blockquote>",
	"<pre", "</pre>", "<table", "</table>", "<tr", "</tr>", "<td", "</td>", "<th", "</th>"}

// highlightInsert wraps text content in mark tags while preserving HTML structure
func highlightInsert(text string) string {
	// If the text doesn't contain any HTML tags, just wrap it
	if !strings.Contains(text, "<") {
		return `<mark class="revision-insert">` + text + `</mark>`
	}

	// Check if the insertion contains block-level tags
	// If so, we shouldn't wrap the whole thing as it would break HTML structure
	for _, tag := range blockLevelTags {
		if strings.Contains(text, tag) {
			// Contains block-level tags - just return as-is without highlighting
			// to avoid creating invalid HTML
			return text
		}
	}

	// Only contains inline tags - safe to wrap the whole thing
	return `<mark class="revision-insert">` + text + `</mark>`
}
