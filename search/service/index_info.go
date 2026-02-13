package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"search/config"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

type IndexInfo struct {
	Slug         string     `json:"slug"`
	Name         string     `json:"name"`
	LastModified *time.Time `json:"last_modified"`
	ArchiveDate  *time.Time `json:"archive_date"`
	Content      string     `json:"content"`
}

func getIndexInfoBatch(index int, count int) ([]IndexInfo, error) {
	url := fmt.Sprintf("%s/indexable-pages?index=%d&count=%d", config.WikiURL, index, count)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var indexInfo []IndexInfo
	err = json.Unmarshal(body, &indexInfo)
	if err != nil {
		return nil, err
	}

	return indexInfo, nil
}

func buildIndexMapping() mapping.IndexMapping {
	docMapping := bleve.NewDocumentMapping()

	nameMapping := bleve.NewTextFieldMapping()
	nameMapping.Analyzer = "en"
	nameMapping.Store = false
	nameMapping.Index = true
	docMapping.AddFieldMappingsAt("name", nameMapping)

	contentMapping := bleve.NewTextFieldMapping()
	contentMapping.Analyzer = "en"
	contentMapping.Store = false
	contentMapping.Index = true
	docMapping.AddFieldMappingsAt("content", contentMapping)

	slugMapping := bleve.NewTextFieldMapping()
	slugMapping.Analyzer = "keyword"
	slugMapping.Store = false
	slugMapping.Index = true
	docMapping.AddFieldMappingsAt("slug", slugMapping)

	dateMapping := bleve.NewDateTimeFieldMapping()
	dateMapping.Store = false
	dateMapping.Index = true
	dateMapping.IncludeInAll = false
	docMapping.AddFieldMappingsAt("last_modified", dateMapping)
	docMapping.AddFieldMappingsAt("archive_date", dateMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = docMapping
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping
}
