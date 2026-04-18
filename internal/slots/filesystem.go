package slots

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name    string
	Path    string
	Content string
}

type FilesystemLoader struct {
	root string
}

func NewFilesystemLoader(root string) *FilesystemLoader {
	return &FilesystemLoader{root: root}
}

func (l *FilesystemLoader) LoadSpecialistSlots(specialistID string) ([]File, error) {
	dir := filepath.Join(l.root, specialistID)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read specialist slot dir: %w", err)
	}

	files := make([]File, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			continue
		}
		path := filepath.Join(dir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read slot file %s: %w", name, err)
		}
		files = append(files, File{Name: name, Path: path, Content: string(content)})
	}

	return files, nil
}

func (l *FilesystemLoader) Exists(specialistID string) bool {
	_, err := os.Stat(filepath.Join(l.root, specialistID))
	return err == nil || !os.IsNotExist(err)
}
