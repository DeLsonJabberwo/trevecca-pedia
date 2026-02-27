package handlers

import (
	"context"
	"io"
	"net/http"
	"time"
	wikierrors "wiki/errors"
	"wiki/requests"
	"wiki/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewPageHandler(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		werr, is := wikierrors.AsWikiError(err)
		if !is {
			werr = wikierrors.InternalError(err)
		}
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}
	defer db.Close()
	dataDir := utils.GetDataDir()

	var newPageReq utils.NewPageRequest
	err = c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	file, err := c.FormFile("new_page")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	f, err := file.Open()
	if err != nil {
		werr := wikierrors.FilesystemError(err)
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}
	defer f.Close()
	newPageBytes, err := io.ReadAll(f)
	if err != nil {
		werr := wikierrors.InternalError(err)
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}
	newPageReq.Slug = c.PostForm("slug")
	newPageReq.Name = c.PostForm("name")
	newPageReq.Author = c.PostForm("author")

	archiveDateStr := c.PostForm("archive_date")
	if archiveDateStr != "" {
		archiveDate, err := time.Parse("2006-01-02", archiveDateStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "bad request format",
			})
			return
		}
		newPageReq.ArchiveDate = &archiveDate
	}

	newPageReq.Content = string(newPageBytes)

	err = utils.CreateNewPage(ctx, db, dataDir, newPageReq)
	if err != nil {
		werr, is := wikierrors.AsWikiError(err)
		if !is {
			werr = wikierrors.InternalError(err)
		}
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}
	c.Status(http.StatusOK)
}

func DeletePageHandler(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		werr, is := wikierrors.AsWikiError(err)
		if !is {
			werr = wikierrors.InternalError(err)
		}
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}
	defer db.Close()
	dataDir := utils.GetDataDir()

	var delReq utils.DeletePageRequest
	err = c.Request.ParseForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}

	delReq.Slug = c.PostForm("slug")
	delReq.User = c.PostForm("user")

	err = requests.DeletePage(ctx, db, dataDir, delReq)
	if err != nil {
		werr, is := wikierrors.AsWikiError(err)
		if !is {
			werr = wikierrors.InternalError(err)
		}
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}

	c.Status(http.StatusOK)
}

func NewRevisionHandler(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		werr, is := wikierrors.AsWikiError(err)
		if !is {
			werr = wikierrors.InternalError(err)
		}
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}
	defer db.Close()
	dataDir := utils.GetDataDir()

	var revReq utils.RevisionRequest
	err = c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	file, err := c.FormFile("new_page")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	f, err := file.Open()
	if err != nil {
		werr := wikierrors.FilesystemError(err)
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}
	defer f.Close()
	newPageBytes, err := io.ReadAll(f)
	if err != nil {
		werr := wikierrors.InternalError(err)
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}
	revReq.PageId = c.PostForm("page_id")
	revReq.Author = c.PostForm("author")
	revReq.NewPage = string(newPageBytes)

	err = requests.PostRevision(ctx, db, dataDir, revReq)
	if err != nil {
		werr, is := wikierrors.AsWikiError(err)
		if !is {
			werr = wikierrors.InternalError(err)
		}
		c.AbortWithStatusJSON(werr.Code, gin.H{
			"error": werr.Details,
		})
		return
	}

	c.Status(http.StatusOK)
}

func EditPageHandler(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		werr, is := wikierrors.AsWikiError(err)
		if !is {
			werr = wikierrors.InternalError(err)
		}
		c.AbortWithStatusJSON(werr.Code, gin.H{"error": werr.Details})
		return
	}
	defer db.Close()
	dataDir := utils.GetDataDir()

	pageId := c.Param("id")
	if _, err := uuid.Parse(pageId); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid page id",
		})
		return
	}

	content := c.PostForm("content")
	if content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "content cannot be empty",
		})
		return
	}

	revReq := utils.RevisionRequest{
		PageId:  pageId,
		Author:  c.PostForm("author"),
		NewPage: content,
	}

	err = requests.PostRevision(ctx, db, dataDir, revReq)
	if err != nil {
		werr, is := wikierrors.AsWikiError(err)
		if !is {
			werr = wikierrors.InternalError(err)
		}
		c.AbortWithStatusJSON(werr.Code, gin.H{"error": werr.Details})
		return
	}

	c.Status(http.StatusOK)
}

