package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func GetPageContent(dataDir string, pageId uuid.UUID) (string, error) {
	content, err := os.ReadFile(filepath.Join(dataDir, "pages", fmt.Sprintf("%s.md", pageId)))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetRevisionContent(dataDir string, revId uuid.UUID) (string, error) {
	if revId == uuid.Nil {
		return "", nil
	}
	content, err := os.ReadFile(filepath.Join(dataDir, "revisions", fmt.Sprintf("%s.txt", revId)))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetSnapshotContent(dataDir string, snapId uuid.UUID) (string, error) {
	content, err := os.ReadFile(filepath.Join(dataDir, "snapshots", fmt.Sprintf("%s.md", snapId)))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
