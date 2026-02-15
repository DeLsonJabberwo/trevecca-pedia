package service

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

type SearchService struct {
	index bleve.Index
}

func NewSearchService(indexPath string) (*SearchService, error) {
	idx, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := buildIndexMapping()
		idx, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, fmt.Errorf("failed to create index: %w", err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open index: %w", err)
	}

	return &SearchService{index: idx}, nil
}

type PageDocument struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *SearchService) IndexAll(pagesDir string) error {
	batch := s.index.NewBatch()
	batchLen := 10
	currInd := 0

	for {
		indexInfoBatch, err := getIndexInfoBatch(currInd, batchLen)
		if err != nil {
			return err
		}

		for _, indexInfo := range indexInfoBatch {
			batch.Index(indexInfo.Slug, indexInfo)
		}

		if len(indexInfoBatch) < batchLen {
			break
		}

		currInd += batchLen
	}

	return s.index.Batch(batch)
}

func (s *SearchService) Search(queryString string, from, size int) (*bleve.SearchResult, error) {
	nameQuery := bleve.NewMatchQuery(queryString)
	nameQuery.SetField("name")
	nameQuery.SetBoost(5.0)

	contentQuery := bleve.NewMatchQuery(queryString)
	contentQuery.SetField("content")
	contentQuery.SetBoost(1.0)

	slugQuery := bleve.NewMatchQuery(queryString)
	slugQuery.SetField("slug")
	slugQuery.SetBoost(2.0)

	disjunctionQuery := bleve.NewDisjunctionQuery(nameQuery, contentQuery, slugQuery)

	req := bleve.NewSearchRequest(disjunctionQuery)
	req.From = from
	req.Size = size

	return s.index.Search(req)
}
