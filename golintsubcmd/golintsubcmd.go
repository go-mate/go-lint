package golintsubcmd

import (
	"time"

	"github.com/go-mate/go-lint/goworkcilint"
	"github.com/go-mate/go-work/workcfg"
	"github.com/spf13/cobra"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func NewRunCmd(config *workcfg.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "golangci-lint run",
		Long:  "golangci-lint run",
		Run: func(cmd *cobra.Command, args []string) {
			Run(config, time.Minute*5)
		},
	}
}

func Run(config *workcfg.WorksExec, timeout time.Duration) {
	projects := config.Subprojects()
	zaplog.SUG.Debugln("lint", neatjsons.S(projects))

	output := rese.V1(config.GetNewCommand().Exec("golangci-lint", "version"))
	zaplog.SUG.Debugln(string(output))

	result := goworkcilint.Run(config.GetNewCommand(), projects, timeout)
	goworkcilint.DebugIssues(projects, result)
}
