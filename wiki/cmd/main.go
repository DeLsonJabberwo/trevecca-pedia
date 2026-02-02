package main

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"wiki/database"
	wikierrors "wiki/errors"
	"wiki/requests"
	"wiki/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	ctx, db, dataDir := setup()
	defer db.Close()

	// /pages?category={cat}&index={ind}&num={num}
	r.GET("/pages", func(c *gin.Context) {
		catQuery := c.DefaultQuery("category", "")
		cat := database.ValidateCategory(ctx, db, catQuery)
		ind, err := strconv.Atoi(c.DefaultQuery("index", "0"))
		if err != nil {
			ind = 0
		}
		num, err := strconv.Atoi(c.DefaultQuery("num", "10"))
		if err != nil {
			num = 10
		}
		var pages []database.PageInfo
		if cat == 0 {
			pages, err = requests.GetPages(ctx, db, ind, num)
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
			pages, err = requests.GetPagesCategory(ctx, db, cat, ind, num)
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
	})

	r.GET("/pages/:id", func(c *gin.Context) {
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
	})

	r.GET("/pages/:id/revisions", func(c *gin.Context) {
		pageId := c.Param("id")
		ind, err := strconv.Atoi(c.DefaultQuery("index", "0"))
		if err != nil {
			ind = 0
		}
		num, err := strconv.Atoi(c.DefaultQuery("num", "10"))
		if err != nil {
			num = 10
		}
		revisions, err := requests.GetRevisions(ctx, db, pageId, ind, num)
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
	})

	r.GET("/pages/:id/revisions/:rev", func(c *gin.Context) {
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
	})

	// POST

	r.POST("/pages/new", func(c *gin.Context) {
		var newPageReq utils.NewPageRequest
		err := c.Request.ParseMultipartForm(32 << 20)
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

		// Handle optional archive_date
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
	})

	r.POST("/pages/:id/delete", func(c *gin.Context) {
		var delReq utils.DeletePageRequest
		err := c.Request.ParseMultipartForm(32 << 20)
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

	})

	r.POST("/pages/:id/revisions", func(c *gin.Context) {
		var revReq utils.RevisionRequest
		err := c.Request.ParseMultipartForm(32 << 20)
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
	})

	r.Run(":9454")
}

func setup() (context.Context, *sql.DB, string) {
	ctx := context.Background()

	var connStr = "host=localhost port=5432 dbname=wiki user=wiki_user password=myatt sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	dataDir := filepath.Join("..", "wiki-fs")

	return ctx, db, dataDir
}
