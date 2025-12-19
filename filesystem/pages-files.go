package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)


func GetPage(dataDir string, pageId uuid.UUID) (string, error) {
	file, err := os.ReadFile(filepath.Join(dataDir, "pages", fmt.Sprintf("%s.md", pageId)))
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func GetRevision(dataDir string, revId uuid.UUID) (string, error) {
	file, err := os.ReadFile(filepath.Join(dataDir, "revisions", fmt.Sprintf("%s.txt", revId)))
	if err != nil {
		return "", err
	}
	return string(file), nil
}

