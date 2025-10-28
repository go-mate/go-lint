// Package main: go-lint CLI tool for intelligent golangci-lint execution
// Provides smart linting for Go projects with workspace support and debug modes
// Supports single package, multiple roots, and workspace-wide linting operations
//
// main: go-lint CLI 工具，用于智能化 golangci-lint 执行
// 为 Go 项目提供智能 linting，支持工作区和调试模式
// 支持单包、多根路径和工作区范围的 linting 操作
package main

import (
	"os"
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/go-mate/go-work/workspath"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

// Usage examples:
// go-lint run
// go-lint run --debug
// go-lint run --debug=true
// go-lint run --debug=1
// go-lint run --debug=0
// go-lint
//
// 使用示例:
// go-lint run
// go-lint run --debug
// go-lint run --debug=true
// go-lint run --debug=1
// go-lint run --debug=0
// go-lint
func main() {
	// Get current working DIR as the base path for linting operations
	// 获取当前工作 DIR 作为 linting 操作的基础路径
	workPath := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(workPath))

	// Create root command with basic linting functionality
	// 创建具有基础 linting 功能的根命令
	rootCmd := cobra.Command{
		Use:   "go-lint",
		Short: "Smart Go linting tool with workspace support",
		Long:  "go-lint provides intelligent golangci-lint execution with multi-project workspace support and debug modes",
		Run: func(cmd *cobra.Command, args []string) {
			run(workPath, false)
		},
	}
	rootCmd.AddCommand(newRunCmd(workPath))
	must.Done(rootCmd.Execute())
}

// newRunCmd creates the 'run' subcommand with debug mode support
// Provides detailed control over linting execution with verbose output options
//
// newRunCmd 创建带有调试模式支持的 'run' 子命令
// 提供对 linting 执行的详细控制和详细输出选项
func newRunCmd(workPath string) *cobra.Command {
	var debugMode bool
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Execute linting with optional debug mode",
		Long:  "Run golangci-lint on the current package and submodules with optional debug output for detailed analysis",
		Run: func(cmd *cobra.Command, args []string) {
			run(workPath, debugMode)
		},
	}
	cmd.Flags().BoolVarP(&debugMode, "debug", "", false, "enable debug mode for verbose output and detailed analysis")
	return cmd
}

// run executes the main linting logic with workspace discovery and configuration
// Focuses on current package and submodules, excluding non-Go project-wide linting
//
// run 执行主要的 linting 逻辑，包含工作区发现和配置
// 专注于当前包和子模块，排除非 Go 项目级别的 linting
func run(workPath string, debugMode bool) {
	// Set global debug mode for all linting operations
	// 为所有 linting 操作设置全局调试模式
	golint.SetDebugMode(debugMode)

	// Configure workspace options for targeted linting
	// 配置工作区选项进行目标 linting
	config := workspath.NewOptions().
		WithIncludeCurrentProject(false).
		WithIncludeCurrentPackage(true).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(debugMode)

	// Discover module paths based on workspace configuration
	// 基于工作区配置发现模块路径
	moduleRoots := workspath.GetModulePaths(workPath, config)

	// Execute linting on discovered modules with timeout protection
	// 对发现的模块执行 linting，带有超时保护
	result := golint.RootsRun(moduleRoots, time.Minute*5)

	// Display detailed results and issues for analysis
	// 显示详细结果和问题进行分析
	result.DebugIssues()
}
