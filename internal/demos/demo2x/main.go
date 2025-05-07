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
// go run main.go run
func main() {
	projectPath := runpath.PARENT.Up(3)
	zaplog.SUG.Debugln(projectPath)

	workspace := workcfg.NewWorkspace("", []string{projectPath})

	execConfig := osexec.NewCommandConfig().WithBash().WithDebugMode(true)
	workspaces := []*workcfg.Workspace{workspace}

	config := workcfg.NewWorksExec(execConfig, workspaces)

	// 定义根命令
	var rootCmd = &cobra.Command{
		Use:   "lint", // 根命令的名称
		Short: "lint",
		Long:  "lint",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("lint")
		},
	}
	rootCmd.AddCommand(golintsubcmd.NewRunCmd(config))

	must.Done(rootCmd.Execute())
}
