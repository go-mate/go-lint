package golintworks_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/golintworks"
	"github.com/go-mate/go-work/workcfg"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestRun(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := workcfg.NewWorkspace("", []string{projectPath})

	execConfig := osexec.NewCommandConfig().WithDebugMode(true)
	workspaces := []*workcfg.Workspace{
		workspace,
	}
	config := workcfg.NewWorksExec(execConfig, workspaces)

	resMap := golintworks.Run(config, time.Minute)
	t.Log(neatjsons.S(resMap))

	golintworks.DebugIssues(config, resMap)
}
