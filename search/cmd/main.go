package main

import (
	"log"
	"search/config"
	"search/handlers"
	"search/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	s, err := service.NewSearchService(config.IndexDir)
	if err != nil {
		log.Fatalf("Couldn't create search service: %s\n", err)
	}

	err = s.IndexAll()
	if err != nil {
		log.Fatalf("Couldn't index search service: %s\n", err)
	}

	handlers.SetSearchService(s)

	r.GET("/search", handlers.SearchHandler)
	r.POST("/reindex", handlers.ReindexHandler)

	r.Run(":7724")
}
