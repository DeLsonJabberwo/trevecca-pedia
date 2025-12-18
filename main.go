package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"wiki/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	log.Printf("testGetPageInfo(ctx, db):\n")
	testGetPageInfo(ctx, db)
	fmt.Println()
	log.Printf("testGetPageNameUUID(ctx, db):\n")
	testGetPageNameUUID(ctx, db)
	fmt.Println()
	log.Printf("testGetPageRevisionsInfo(ctx, db):\n")
	testGetPageRevisionsInfo(ctx, db)

	fmt.Println()
	r.Run(":8080")
}

func testGetPageInfo(ctx context.Context, db *sql.DB) {
	pageUUID, err := uuid.Parse("07918316-875e-4581-87ab-5b8d1d8bdd3a")
	if err != nil {
		log.Fatal(err)
	}
	testPage, err := database.GetPageInfo(ctx, db, pageUUID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v\n", testPage)
}

func testGetPageNameUUID(ctx context.Context, db *sql.DB) {
	res, err := database.GetPageNameUUIDs(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Name\t\tUUID\n")
	for _, r := range res {
		log.Printf("%s\t\t%s\n", r.Name, r.UUID)
	}
}

func testGetPageRevisionsInfo(ctx context.Context, db *sql.DB) {
	pageUUID, err := uuid.Parse("07918316-875e-4581-87ab-5b8d1d8bdd3a")
	if err != nil {
		log.Fatal(err)
	}
	pageRevs, err := database.GetPageRevisionsInfo(ctx, db, pageUUID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("{UUID,\tDateTime,\tAuthor}\n")
	for _, r := range pageRevs {
		log.Printf("{%s,\t%s,\t%s}\n", r.UUID, r.DateTime, r.Author)
	}
}

