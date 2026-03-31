package main

import (
	"moderation/config"
	"moderation/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	r.GET("/flagged-users", handlers.GetFlaggedUsers)
	r.GET("/silenced-users", handlers.GetSilencedUsers)

	port := config.GetEnv("MOD_SERVICE_PORT", "6633")
	r.Run(":" + port)
}
