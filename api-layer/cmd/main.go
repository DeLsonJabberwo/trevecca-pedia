package main

import (
	"api-layer/handlers/wiki"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	r.GET("/v1/wiki/pages", wiki.GetPages)
	r.GET("/v1/wiki/pages/:id", wiki.GetPage)
	r.GET("/v1/wiki/pages/:id/revisions", wiki.GetPageRevisions)
	r.GET("/v1/wiki/pages/:id/revisions/:rev", wiki.GetPageRevision)

	r.POST("/v1/wiki/pages/new", wiki.PostNewPage)
	r.POST("/v1/wiki/pages/:id/revisions", wiki.PostPageRevision)


	r.Run(":2745")
}
