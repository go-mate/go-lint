package golint

import (
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/go-mate/go-lint/golintroots"
	"github.com/go-mate/go-lint/golintworks"
	"github.com/go-mate/go-work/workspace"
	"github.com/yyle88/osexec"
)

func Run(path string, timeout time.Duration) *golangcilint.Result {
	execConfig := osexec.NewExecConfig().WithDebugMode(osexec.NewDebugMode(debugModeOpen))

	result := golangcilint.Run(execConfig, path, timeout)
	if debugModeOpen {
		result.DebugIssues()
	}
	return result
}

func RootsRun(roots []string, timeout time.Duration) *golintroots.Result {
	execConfig := osexec.NewExecConfig().WithDebugMode(osexec.NewDebugMode(debugModeOpen))

	result := golintroots.Run(execConfig, roots, timeout)
	if debugModeOpen {
		result.DebugIssues()
	}
	return result
}

func WorksRun(workspaces []*workspace.Workspace, timeout time.Duration) *golintworks.Result {
	execConfig := osexec.NewExecConfig().WithDebugMode(osexec.NewDebugMode(debugModeOpen))

	result := golintworks.Run(execConfig, workspaces, timeout)
	if debugModeOpen {
		result.DebugIssues()
	}
	return result
}
