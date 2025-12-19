package main

import (
	"fmt"
	"log"
	"wiki/filesystem"

	"github.com/google/uuid"
)

func testGetPage(dataDir string) {
	pageId, err := uuid.Parse("07918316-875e-4581-87ab-5b8d1d8bdd3a")
	if err != nil {
		log.Fatal(err)
	}
	page, err := filesystem.GetPage(dataDir, pageId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(page)
}
