package main

import (
	"path/filepath"

	"github.com/go-mate/go-lint/golintworks"
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// go run main.go help
// go run main.go run
func main() {
	projectPath := runpath.PARENT.Up(3)
	osmustexist.MustFile(filepath.Join(projectPath, "go.mod"))
	zaplog.SUG.Debugln(projectPath)

	wsp := workspace.NewWorkspace("", []string{projectPath})

	execConfig := osexec.NewCommandConfig().WithBash().WithDebug()
	workspaces := []*workspace.Workspace{wsp}

	config := worksexec.NewWorksExec(execConfig, workspaces)

	// 定义根命令
	var rootCmd = &cobra.Command{
		Use:   "lint", // 根命令的名称
		Short: "lint",
		Long:  "lint",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("lint")
		},
	}
	rootCmd.AddCommand(golintworks.NewCmd(config))

	must.Done(rootCmd.Execute())
}
