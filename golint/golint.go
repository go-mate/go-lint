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
	golangcilint.DebugIssues(path, result)
	return result
}

func RootsRun(execConfig *osexec.ExecConfig, roots []string, timeout time.Duration) map[string]*golintroots.Result {
	resMap := golintroots.Run(execConfig, roots, timeout)
	golintroots.DebugIssues(roots, resMap)
	return resMap
}

func WorksRun(config *workcfg.WorksExec, timeout time.Duration) map[string]*golintroots.Result {
	resMap := golintworks.Run(config, timeout)
	golintworks.DebugIssues(config, resMap)
	return resMap
}
