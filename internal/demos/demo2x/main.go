package main

import (
	"github.com/go-mate/go-lint/golintsubcmd"
	"github.com/go-mate/go-work/workcfg"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// go run main.go help
// go run main.go lint
func main() {
	projectPath := runpath.PARENT.Up(3)
	zaplog.SUG.Debugln(projectPath)

	workspace := workcfg.NewWorkspace("", []string{projectPath})

	commandConfig := osexec.NewCommandConfig()
	commandConfig.WithBash()
	commandConfig.WithDebugMode(true)

	config := workcfg.NewWorksExec([]*workcfg.Workspace{workspace}, commandConfig)

	// 定义根命令
	var rootCmd = &cobra.Command{
		Use:   "go", // 根命令的名称
		Short: "run",
		Long:  "run",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("run")
		},
	}
	rootCmd.AddCommand(golintsubcmd.NewLintCmd(config))

	must.Done(rootCmd.Execute())
}
