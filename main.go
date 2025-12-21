package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"wiki/requests"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	ctx, db, dataDir := setup()

	r.GET("/pages/:id", func(c *gin.Context) {
		pageId := c.Param("id")
		page, err := requests.GetPage(ctx, db, dataDir, pageId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		c.JSON(http.StatusOK, page)
	})
	
	etcTesting(ctx, db, dataDir)

	r.Run(":8080")
}

func setup() (context.Context, *sql.DB, string) {
	ctx := context.Background()

	var connStr = "host=localhost port=5432 dbname=wiki user=wiki_user password=myatt sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dataDir := filepath.Join(home, "trevecca", "trevecca-pedia", "wiki-fs")

	return ctx, db, dataDir
}

func etcTesting(ctx context.Context, db *sql.DB, dataDir string) {
	fmt.Println()
	fmt.Println()
	log.Printf("Etc. Testing...\n")
	fmt.Println()
	testConnection(ctx, db)

	fmt.Println()
	log.Printf("Testing File System...\n")


	fmt.Println()
	log.Printf("testGetPage(dataDir)\n")
	testGetPage(dataDir)

	fmt.Println()
	fmt.Println()
}

