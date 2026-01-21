package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"wiki/database"
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
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			}
		} else {
			pages, err = requests.GetPagesCategory(ctx, db, cat, ind, num)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
			}
		}
		c.JSON(http.StatusOK, pages)
	})

	r.GET("/pages/:id", func(c *gin.Context) {
		pageId := c.Param("id")
		page, err := requests.GetPage(ctx, db, dataDir, pageId)
		if err != nil && err.Error() == strconv.Itoa(http.StatusNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
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
		if err != nil && err.Error() == strconv.Itoa(http.StatusNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
										"error": err.Error(),
									})
			return
		}
		c.JSON(http.StatusOK, revisions)
	})

	r.GET("/pages/:id/revisions/:rev", func(c *gin.Context) {
		revId := c.Param("rev")
		revision, err := requests.GetRevision(ctx, db, dataDir, revId)
		if err != nil && err.Error() == "404" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
										"error": err.Error(),
									})
		}
		c.JSON(http.StatusOK, revision)
	})

	// POST

	r.POST("/pages/:id/revisions", func(c *gin.Context) {
		var revReq utils.RevisionRequest
		err := c.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		file, err := c.FormFile("new_page")
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		f, _ := file.Open()
		defer f.Close()
		newPageBytes, _ := io.ReadAll(f)
		revReq.PageId = c.PostForm("page_id")
		revReq.Author = c.PostForm("author")
		revReq.NewPage = string(newPageBytes)

		err = requests.PostRevision(ctx, db, dataDir, revReq)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
										"error": err.Error(),
									})
		}
	})
	
	//etcTesting(db, dataDir)

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

func etcTesting(db *sql.DB, dataDir string) {
	fmt.Println()
	fmt.Println()
	log.Printf("Etc. Testing...\n")
	fmt.Println()
	testConnection(db)

	fmt.Println()
	log.Printf("Testing File System...\n")


	fmt.Println()
	log.Printf("testGetPage(dataDir)\n")
	testGetPage(dataDir)

	fmt.Println()
	fmt.Println()
}

