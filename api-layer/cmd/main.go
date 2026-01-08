package main

import (
	"api-layer/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	r.GET("/v1/wiki/pages", handlers.GetPages)
	r.GET("/v1/wiki/pages/:id", handlers.GetPage)
	r.GET("/v1/wiki/pages/:id/revisions", handlers.GetPageRevisions)
	r.GET("/v1/wiki/pages/:id/revisions/:rev", handlers.GetPageRevision)


	r.Run(":2745")
}
