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
	page, err := filesystem.GetPageContent(dataDir, pageId)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(page)
	if page != "" {
		fmt.Println("GetPage(): success")
	}
}
