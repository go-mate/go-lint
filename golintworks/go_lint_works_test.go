package golintworks_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/golintworks"
	"github.com/go-mate/go-work/workcfg"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestRun(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	execConfig := osexec.NewCommandConfig().WithDebugMode(true)
	workspaces := []*workcfg.Workspace{
		workcfg.NewWorkspace("", []string{projectPath}),
	}
	result := golintworks.Run(execConfig, workspaces, time.Minute)
	t.Log(neatjsons.S(result))

	result.DebugIssues()
	require.True(t, result.Success())
	require.Equal(t, result.GetMap().Size(), len(workspaces))
}
