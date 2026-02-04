package main

import (
	"wiki/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	// GET

	// /pages?category={cat}&index={ind}&count={count}
	r.GET("/pages", handlers.PagesHandler)

	r.GET("/pages/:id", handlers.PageHandler)

	r.GET("/pages/:id/revisions", handlers.PageRevisionsHandler)

	r.GET("/pages/:id/revisions/:rev", handlers.PageRevisionHandler)


	// POST

	r.POST("/pages/new", handlers.NewPageHandler)

	r.POST("/pages/:id/delete", handlers.DeletePageHandler)

	r.POST("/pages/:id/revisions", handlers.NewRevisionHandler)

	r.Run(":9454")
}

