package codereview

import (
	"os"
	"path/filepath"
	"strings"
)

func walkDir(dirPath string, callback func(string)) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if isAnalyzedFile(path) {
			callback(path)
		}
		return nil
	})
}

func isAnalyzedFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	supported := map[string]bool{
		".go": true, ".py": true, ".js": true, ".ts": true,
		".rs": true, ".java": true, ".rb": true, ".php": true,
		".c": true, ".h": true, ".cpp": true, ".hpp": true,
		".sh": true, ".yml": true, ".yaml": true, ".json": true,
		".xml": true, ".toml": true, ".ini": true, ".cfg": true,
	}
	return supported[ext]
}