package golint

import (
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/go-mate/go-lint/golintroots"
	"github.com/go-mate/go-lint/golintworks"
	"github.com/go-mate/go-work/workcfg"
	"github.com/yyle88/osexec"
)

func Run(execConfig *osexec.ExecConfig, path string, timeout time.Duration) *golangcilint.Result {
	result := golangcilint.Run(execConfig, path, timeout)
	result.DebugIssues()
	return result
}

func RootsRun(execConfig *osexec.ExecConfig, roots []string, timeout time.Duration) *golintroots.Result {
	result := golintroots.Run(execConfig, roots, timeout)
	result.DebugIssues()
	return result
}

func WorksRun(execConfig *osexec.ExecConfig, workspaces []*workcfg.Workspace, timeout time.Duration) *golintworks.Result {
	result := golintworks.Run(execConfig, workspaces, timeout)
	result.DebugIssues()
	return result
}
