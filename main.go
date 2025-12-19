package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"wiki/database"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	
	var connStr = "host=localhost port=5432 dbname=wiki user=wiki_user password=myatt sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	fmt.Println()
	database.TestConnection(ctx, db)

	
	fmt.Println()
	log.Printf("Testing Database...\n")

	fmt.Println()
	log.Printf("testGetPageInfo(ctx, db):\n")
	testGetPageInfo(ctx, db)
	fmt.Println()
	log.Printf("testGetPageNameUUID(ctx, db):\n")
	testGetPageNameUUID(ctx, db)
	fmt.Println()
	log.Printf("testGetPageRevisionsInfo(ctx, db):\n")
	testGetPageRevisionsInfo(ctx, db)

	fmt.Println()
	fmt.Println()
	log.Printf("Testing File System...\n")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dataDir := filepath.Join(home, "trevecca", "trevecca-pedia", "wiki-fs")

	fmt.Println()
	log.Printf("testGetPage(dataDir)\n")
	testGetPage(dataDir)

	fmt.Println()
	r.Run(":8080")
}

