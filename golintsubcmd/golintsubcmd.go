package golintsubcmd

import (
	"time"

	"github.com/go-mate/go-lint/goworkcilint"
	"github.com/go-mate/go-work/workconfig"
	"github.com/spf13/cobra"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func NewLintCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "lint",
		Short: "golangci-lint run",
		Long:  "golangci-lint run",
		Run: func(cmd *cobra.Command, args []string) {
			Run(config, time.Minute*5)
		},
	}
}

func Run(config *workconfig.WorkspacesExecConfig, timeout time.Duration) {
	projects := config.CollectSubprojectPaths()
	zaplog.SUG.Debugln("lint", neatjsons.S(projects))

	output := rese.V1(config.GetSubCommand("").Exec("golangci-lint", "version"))
	zaplog.SUG.Debugln(string(output))

	resMap := goworkcilint.Run(config.GetSubCommand(""), projects, timeout)
	goworkcilint.DebugIssues(projects, resMap)
}
