package utils

import (
	"path/filepath"

	"github.com/yyle88/osexistpath/osomitexist"
)

func GetProjectPath(currentPath string) (string, string, bool) {
	projectPath := currentPath
	shortMiddle := ""
	for !osomitexist.IsFile(filepath.Join(projectPath, "go.mod")) {
		subName := filepath.Base(projectPath)

		prePath := filepath.Dir(projectPath)
		if prePath == projectPath {
			return "", "", false
		}
		projectPath = prePath
		shortMiddle = filepath.Join(subName, shortMiddle)
	}
	return projectPath, shortMiddle, true
}
