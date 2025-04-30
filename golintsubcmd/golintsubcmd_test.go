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

	config := workcfg.NewWorksExec([]*workcfg.Workspace{workspace}, osexec.NewCommandConfig().WithDebugMode(true))

	golintsubcmd.Run(config, time.Minute)
}
