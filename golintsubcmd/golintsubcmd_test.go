package golintsubcmd_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/golintsubcmd"
	"github.com/go-mate/go-work/workconfig"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestRunDebug(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := workconfig.NewWorkspace("", []string{projectPath})
	workspace.MustCheck()

	workspaces := workconfig.NewWorkspaces(workspace)
	workspaces.MustCheck()

	config := workconfig.NewWorkspacesExecConfig(workspaces, osexec.NewCommandConfig().WithDebugMode(true))
	config.MustCheck()

	golintsubcmd.Run(config, time.Minute)
}
