package main

import (
	"api-layer/config"
	"api-layer/handlers/search"
	"api-layer/handlers/wiki"
	"api-layer/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	// Public endpoints - no auth required
	r.GET("/v1/wiki/pages", wiki.GetPages)
	r.GET("/v1/wiki/pages/:id", wiki.GetPage)
	r.GET("/v1/wiki/pages/:id/revisions", wiki.GetPageRevisions)
	r.GET("/v1/wiki/pages/:id/revisions/:rev", wiki.GetPageRevision)
	r.GET("/v1/wiki/indexable-pages", wiki.GetIndexablePages)

	// Protected endpoints - require valid token and contributor role
	protected := r.Group("/v1/wiki")
	protected.Use(middleware.AuthMiddleware(), middleware.RequireRole("contributor"))
	{
		protected.POST("/pages/new", wiki.PostNewPage)
		protected.POST("/pages/:id/delete", wiki.PostDeletePage)
		protected.POST("/pages/:id/revisions", wiki.PostPageRevision)
	}

	r.GET("/v1/search/search", search.SearchRequest)

	port := config.GetEnv("API_LAYER_PORT", "2745")
	r.Run(fmt.Sprintf(":%s", port))
}
