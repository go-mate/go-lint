package main

import (
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// go run main.go help
// go run main.go run
//
// Demo showing RootsRun usage with cobra command framework
// 演示 RootsRun 的使用，配合 cobra 命令框架
func main() {
	projectPath := runpath.PARENT.Path()
	zaplog.SUG.Debugln(projectPath)

	projectPaths := []string{projectPath}

	// Create root command with usage hint
	// 创建根命令并提供使用提示
	var rootCmd = &cobra.Command{
		Use:   "lint",
		Short: "lint",
		Long:  "lint",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("Use 'lint run' to execute linting")
		},
	}

	// Add run subcommand that executes RootsRun
	// 添加执行 RootsRun 的 run 子命令
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "golangci-lint run",
		Long:  "golangci-lint run",
		Run: func(cmd *cobra.Command, args []string) {
			result := golint.RootsRun(projectPaths, time.Minute*5)
			result.DebugIssues()
		},
	})

	must.Done(rootCmd.Execute())
}
