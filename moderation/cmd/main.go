package main

import (
	"moderation/config"
	"moderation/handlers"
	"moderation/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Validate required secrets at startup - fail fast before accepting traffic
	_ = config.GetJWTSecret()

	r := gin.Default()
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.DebugMode)

	mod := r.Group("")
	mod.Use(middleware.AuthMiddleware(), middleware.RequireRole("moderator"))
	{
		mod.GET("/flagged-users", handlers.GetFlaggedUsers)
		mod.GET("/suspended-users", handlers.GetSuspendedUsers)
		mod.GET("/statuses", handlers.GetStatuses)

		mod.POST("/flag-user", handlers.PostFlagUser)
		mod.POST("/suspend-user", handlers.PostSuspendUser)
		mod.POST("/unflag-user", handlers.PostUnflagUser)
		mod.POST("/unsuspend-user", handlers.PostUnsuspendUser)
	}

	port := config.GetEnv("MOD_SERVICE_PORT", "6633")
	r.Run(":" + port)
}
