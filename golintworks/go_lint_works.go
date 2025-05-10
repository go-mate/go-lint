package golintworks

import (
	"time"

	"github.com/go-mate/go-lint/golintroots"
	"github.com/go-mate/go-work/workcfg"
	"github.com/spf13/cobra"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func NewCmd(config *workcfg.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "golangci-lint run",
		Long:  "golangci-lint run",
		Run: func(cmd *cobra.Command, args []string) {
			DebugIssues(config, Run(config, time.Minute*5))
		},
	}
}

func Run(config *workcfg.WorksExec, timeout time.Duration) map[string]*golintroots.Result {
	projects := config.Subprojects()
	zaplog.SUG.Debugln("golangci-lint run", "PROJECTS", neatjsons.S(projects))

	output := rese.V1(config.GetNewCommand().Exec("golangci-lint", "version"))
	zaplog.SUG.Debugln(string(output))

	return golintroots.Run(config.GetNewCommand(), projects, timeout)
}

func DebugIssues(config *workcfg.WorksExec, resMap map[string]*golintroots.Result) {
	golintroots.DebugIssues(config.Subprojects(), resMap)
}
