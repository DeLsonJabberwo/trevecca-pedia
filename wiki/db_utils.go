package main

import (
	"context"
	"database/sql"
	"log"
	"wiki/database"

	"github.com/google/uuid"
)

func testConnection(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection established.")
}

func testGetPageInfo(ctx context.Context, db *sql.DB) {
	//pageUUID, err := uuid.Parse("07918316-875e-4581-87ab-5b8d1d8bdd3a")
	pageUUID, err := uuid.Parse("60b6b10c-db33-4b4c-9dcf-566f5b3c59a4")
	if err != nil {
		log.Fatal(err)
	}
	testPage, err := database.GetPageInfo(ctx, db, pageUUID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", testPage)
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

