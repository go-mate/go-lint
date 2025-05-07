package golintsubcmd_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/golintsubcmd"
	"github.com/go-mate/go-work/workcfg"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestRunDebug(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := workcfg.NewWorkspace("", []string{projectPath})

	execConfig := osexec.NewCommandConfig().WithDebugMode(true)
	workspaces := []*workcfg.Workspace{
		workspace,
	}
	config := workcfg.NewWorksExec(execConfig, workspaces)

	golintsubcmd.Run(config, time.Minute)
}
