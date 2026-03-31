package handlers

import (
	"context"
	"moderation/auth"
	"moderation/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFlaggedUsers(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "couldn't connect to database",
		})
		return
	}
	ind, err := strconv.Atoi(c.DefaultQuery("ind", "0"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	count, err := strconv.Atoi(c.DefaultQuery("count", "20"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	flaggedUsers, err := auth.ListFlaggedUsers(ctx, db, ind, count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(http.StatusOK, flaggedUsers)
}

func GetSuspendedUsers(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "couldn't connect to database",
		})
		return
	}
	ind, err := strconv.Atoi(c.DefaultQuery("ind", "0"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	count, err := strconv.Atoi(c.DefaultQuery("count", "20"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	silencedUsers, err := auth.ListSuspendedUsers(ctx, db, ind, count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(http.StatusOK, silencedUsers)
}

func GetStatuses(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "couldn't connect to database",
		})
		return
	}
	username := c.DefaultQuery("username", "")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	user, err := auth.GetUser(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}
	statuses := auth.GetUserStatuses(ctx, db, user)
	c.JSON(http.StatusOK, statuses)
}
