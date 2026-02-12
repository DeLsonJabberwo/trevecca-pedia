package main

import (
	"web/config"
	"web/wiki"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	r.Static("./static", "static")
	r.GET("/", wiki.GetHome)
	r.GET("/pages/:id", wiki.GetPage)

	port := config.GetEnv("WEB_SERVICE_PORT", "8080")
	r.Run(":" + port)
}
