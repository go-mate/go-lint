package golintworks_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/golintworks"
	"github.com/go-mate/go-lint/internal/gofrsflock"
	"github.com/go-mate/go-work/workspace"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

func TestRun(t *testing.T) {
	flock := rese.P1(gofrsflock.Lock(gofrsflock.CurrentPath(), 5*time.Minute+30*time.Second))
	defer rese.F0(flock.Unlock)

	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	execConfig := osexec.NewCommandConfig().WithDebug()
	workspaces := []*workspace.Workspace{
		workspace.NewWorkspace("", []string{projectPath}),
	}
	result := golintworks.Run(execConfig, workspaces, time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
	require.Equal(t, result.GetMap().Size(), len(workspaces))
}
