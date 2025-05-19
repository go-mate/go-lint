package golint_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/go-mate/go-work/workspace"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

var projectPath string

func TestMain(m *testing.M) {
	path := runpath.PARENT.Up(1)
	osmustexist.MustFile(filepath.Join(path, "go.mod"))
	zaplog.SUG.Debugln(path)

	projectPath = path
	m.Run()
}

func TestRun(t *testing.T) {
	result := golint.Run(projectPath, 5*time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
}

func TestRootsRun(t *testing.T) {
	result := golint.RootsRun([]string{projectPath}, 5*time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
}

func TestWorksRun(t *testing.T) {
	workspaces := []*workspace.Workspace{
		workspace.NewWorkspace("", []string{projectPath}),
	}
	result := golint.WorksRun(workspaces, 5*time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
}
