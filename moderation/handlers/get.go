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
	}
	ind, err := strconv.Atoi(c.DefaultQuery("ind", "0"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
	}
	count, err := strconv.Atoi(c.DefaultQuery("count", "20"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
	}
	flaggedUsers, err := auth.ListFlaggedUsers(ctx, db, ind, count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
	}
	c.JSON(http.StatusOK, flaggedUsers)
}

func GetSilencedUsers(c *gin.Context) {
	ctx := context.Background()
	db, err := utils.GetDatabase()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "couldn't connect to database",
		})
	}
	ind, err := strconv.Atoi(c.DefaultQuery("ind", "0"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
	}
	count, err := strconv.Atoi(c.DefaultQuery("count", "20"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
	}
	silencedUsers, err := auth.ListSilencedUsers(ctx, db, ind, count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
	}
	c.JSON(http.StatusOK, silencedUsers)
}
