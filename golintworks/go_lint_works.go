package golintworks

import (
	"time"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/go-mate/go-lint/golangcilint"
	"github.com/go-mate/go-lint/golintroots"
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func NewCmd(config *worksexec.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "golangci-lint run",
		Long:  "golangci-lint run",
		Run: func(cmd *cobra.Command, args []string) {
			result := Run(config.GetNewCommand(), config.GetWorkspaces(), time.Minute*5)
			result.DebugIssues()
		},
	}
}

type Result struct {
	resMap *linkedhashmap.Map[string, *golangcilint.Result]
}

func (R *Result) Success() bool {
	for _, res := range R.resMap.Values() {
		if !res.Success() {
			return false
		}
	}
	return true
}

func (R *Result) GetMap() *linkedhashmap.Map[string, *golangcilint.Result] {
	return R.resMap
}

func Run(execConfig *osexec.ExecConfig, workspaces []*workspace.Workspace, timeout time.Duration) *Result {
	zaplog.SUG.Debugln("golangci-lint run", "WORKSPACES", neatjsons.S(workspaces))

	output := rese.V1(execConfig.ShallowClone().Exec("golangci-lint", "version"))
	zaplog.SUG.Debugln(string(output))

	resMap := linkedhashmap.New[string, *golangcilint.Result]()
	for _, workspace := range workspaces {
		result := golintroots.Run(execConfig, workspace.Projects, timeout)
		result.DebugIssues()
		result.GetMap().Each(func(path string, value *golangcilint.Result) {
			resMap.Put(path, value)
		})
	}

	return &Result{resMap: resMap}
}

func (R *Result) DebugIssues() {
	if wrongCount := golintroots.CountIssues(R.resMap); wrongCount == 0 {
		eroticgo.GREEN.ShowMessage("SUCCESS")
		return
	}
	//首先显示详细的
	golintroots.DebugIssues1(R.resMap)

	//接着显示缩略的
	golintroots.DebugIssues2(R.resMap)
}
