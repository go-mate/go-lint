package main

import (
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// go run main.go
//
// Demo showing BatchRun usage with custom configuration and project slice
// 演示 BatchRun 的使用，带有自定义配置和项目切片
func main() {
	projectPath := runpath.PARENT.Path()
	zaplog.SUG.Debugln(projectPath)

	// Execute batch linting with debug config and display results
	// 使用调试配置执行批量 linting 并显示结果
	projectPaths := []string{projectPath}
	result := golint.BatchRun(osexec.NewCommandConfig().WithDebug(), projectPaths, time.Minute*5)
	result.DebugIssues()
}
