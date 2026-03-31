package handlers

import (
	"context"
	"moderation/auth"
	"moderation/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostFlagUser(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "couldn't connect to database",
		})
		return
	}

	username := c.DefaultPostForm("username", "")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
		return
	}

	user, err := auth.GetUser(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
		return
	}

	err = auth.FlagUser(ctx, db, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to flag user"})
		return
	}

	c.Status(http.StatusOK)
}

func PostSuspendUser(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "couldn't connect to database",
		})
		return
	}
	username := c.DefaultPostForm("username", "")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
		return
	}
	user, err := auth.GetUser(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
		return
	}
	err = auth.SuspendUser(ctx, db, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to silence user"})
		return
	}
	c.Status(http.StatusOK)
}

func PostUnflagUser(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "couldn't connect to database",
		})
		return
	}
	username := c.DefaultPostForm("username", "")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
		return
	}
	user, err := auth.GetUser(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
		return
	}
	err = auth.UnFlagUser(ctx, db, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to unflag user"})
		return
	}
	c.Status(http.StatusOK)
}

func PostUnsuspendUser(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "couldn't connect to database",
		})
		return
	}
	username := c.DefaultPostForm("username", "")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
		return
	}
	user, err := auth.GetUser(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
		return
	}
	err = auth.UnSuspendUser(ctx, db, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to unsilence user"})
		return
	}
	c.Status(http.StatusOK)
}
