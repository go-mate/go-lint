package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/runpath"
)

func TestGetProjectPath(t *testing.T) {
	path := runpath.PARENT.Path()
	t.Log(path)
	projectPath, shortMiddle, isGoModule := GetProjectPath(path)
	require.True(t, isGoModule)
	t.Log(projectPath)
	t.Log(shortMiddle)

	require.Equal(t, path, filepath.Join(projectPath, shortMiddle))
}
