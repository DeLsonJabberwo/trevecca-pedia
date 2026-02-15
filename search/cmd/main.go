package main

import (
	"log"
	"os"
	"search/handlers"
	"search/service"

	"github.com/gin-gonic/gin"
)

func main() {
	pagesDir := os.Getenv("PAGES_DIR")
	if pagesDir == "" {
		pagesDir = "../wiki-fs/pages"
	}

	indexDir := os.Getenv("INDEX_DIR")
	if indexDir == "" {
		indexDir = "../wiki-fs/index"
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	s, err := service.NewSearchService(indexDir)
	if err != nil {
		log.Fatalf("Couldn't create search service: %s\n", err)
	}

	err = s.IndexAll(pagesDir)
	if err != nil {
		log.Fatalf("Couldn't index search service: %s\n", err)
	}

	handlers.SetSearchService(s)

	r.GET("/search", handlers.SearchHandler)
	r.POST("/reindex", handlers.ReindexHandler)

	r.Run(":7724")
}
