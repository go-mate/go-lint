package golint_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/go-mate/go-work/workcfg"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
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
	result := golint.Run(osexec.NewExecConfig().WithDebug(), projectPath, 5*time.Minute)
	require.True(t, result.Success())
}

func TestRootsRun(t *testing.T) {
	resMap := golint.RootsRun(osexec.NewExecConfig().WithDebug(), []string{projectPath}, 5*time.Minute)
	require.Len(t, resMap, 0)
}

func TestWorksRun(t *testing.T) {
	worksExec := workcfg.NewWorksExec(osexec.NewExecConfig(), []*workcfg.Workspace{
		workcfg.NewWorkspace("", []string{projectPath}),
	})
	resMap := golint.WorksRun(worksExec, 5*time.Minute)
	require.Len(t, resMap, 0)
}
