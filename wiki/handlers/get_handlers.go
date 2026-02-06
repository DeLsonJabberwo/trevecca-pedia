package handlers

import (
	"context"
	"net/http"
	"strconv"
	"wiki/database"
	wikierrors "wiki/errors"
	"wiki/requests"
	"wiki/utils"

	"github.com/gin-gonic/gin"
)

func PagesHandler(c *gin.Context) {
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

	catQuery := c.DefaultQuery("category", "")
	cat := database.ValidateCategory(ctx, db, catQuery)
	ind, err := strconv.Atoi(c.DefaultQuery("index", "0"))
	if err != nil {
		ind = 0
	}
	count, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		count = 10
	}
	var pages []database.PageInfo
	if cat == 0 {
		pages, err = requests.GetPages(ctx, db, ind, count)
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
	} else {
		pages, err = requests.GetPagesCategory(ctx, db, cat, ind, count)
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
	}
	c.JSON(http.StatusOK, pages)
}

func PageHandler(c *gin.Context) {
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

	defer db.Close()
	pageId := c.Param("id")
	page, err := requests.GetPage(ctx, db, dataDir, pageId)
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
	c.JSON(http.StatusOK, page)
}

func PageRevisionsHandler(c *gin.Context) {
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

	pageId := c.Param("id")
	ind, err := strconv.Atoi(c.DefaultQuery("index", "0"))
	if err != nil {
		ind = 0
	}
	count, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		count = 10
	}
	revisions, err := requests.GetRevisions(ctx, db, pageId, ind, count)
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
	c.JSON(http.StatusOK, revisions)
}

func PageRevisionHandler(c *gin.Context) {
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
	revId := c.Param("rev")
	revision, err := requests.GetRevision(ctx, db, dataDir, revId)
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
	c.JSON(http.StatusOK, revision)
}
