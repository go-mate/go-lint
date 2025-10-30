// Package golint: Unified interface to execute golangci-lint across different scopes
// Provides consistent API with single path, multiple roots, and workspace linting operations
// Handles debug mode management and timeout configuration in each linting situation
//
// golint: golangci-lint 跨不同范围执行的统一接口
// 为单路径、多根路径和工作区 linting 操作提供一致的 API
// 处理所有 linting 场景的调试模式管理和超时配置
package golint

import (
	"fmt"
	"time"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/go-mate/go-lint/golangcilint"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/osexec"
)

// Run executes golangci-lint on a single path with timeout and debug support
// Returns structured results with issue details and execution status
//
// Run 在单个路径上执行 golangci-lint，支持超时和调试
// 返回带有问题详情和执行状态的结构化结果
func Run(path string, timeout time.Duration) *golangcilint.Result {
	execConfig := osexec.NewExecConfig().WithDebugMode(osexec.NewDebugMode(debugModeOpen))

	result := golangcilint.Run(execConfig, path, timeout)
	if debugModeOpen {
		result.DebugIssues()
	}
	return result
}

// WorksRun executes golangci-lint across multiple project paths
// This is an alias to RootsRun with semantic naming for multi-project usage
//
// WorksRun 跨多个项目路径执行 golangci-lint
// 这是 RootsRun 的别名，用于多项目场景的语义化命名
func WorksRun(projectPaths []string, timeout time.Duration) *Result {
	return RootsRun(projectPaths, timeout)
}

// RootsRun executes golangci-lint on multiple paths in batch mode
// Processes multiple modules with consolidated result reporting
//
// RootsRun 在多个路径上批量执行 golangci-lint
// 高效处理多个模块并提供合并的结果报告
func RootsRun(roots []string, timeout time.Duration) *Result {
	execConfig := osexec.NewExecConfig().WithDebugMode(osexec.NewDebugMode(debugModeOpen))

	result := BatchRun(execConfig, roots, timeout)
	if debugModeOpen {
		result.DebugIssues()
	}
	return result
}

// BatchRun executes golangci-lint on multiple paths with custom execution configuration
// Base function that provides flexible execution settings
//
// BatchRun 在多个路径上执行 golangci-lint，支持自定义执行配置
// 底层函数，在执行设置方面提供灵活性
func BatchRun(execConfig *osexec.ExecConfig, roots []string, timeout time.Duration) *Result {
	var resMap = linkedhashmap.New[string, *golangcilint.Result]()
	for idx, path := range roots {
		fmt.Println(eroticgo.BLUE.Sprint("(", idx, "/", len(roots), ")"))
		// Show command before execution with progress
		// 执行前显示命令和进度
		golangcilint.ShowCommandMessage(path)

		result := golangcilint.Run(execConfig, path, timeout)
		resMap.Put(path, result)

		// Show result immediately after each project with progress
		// 每个项目检测后立即显示结果和进度
		result.ShowOutlineMessage()
	}
	return NewResult(resMap)
}
