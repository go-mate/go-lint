// Package golint: Batch linting with aggregated results
// Executes golangci-lint across multiple paths with consolidated reporting
// Maintains execution sequence and provides comprehensive success/failure analysis
//
// golint: 批量 linting，支持聚合结果
// 跨多个路径执行 golangci-lint，提供合并报告
// 维护执行顺序并提供全面的成功/失败分析
package golint

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/go-mate/go-lint/golangcilint"
	"github.com/yyle88/eroticgo"
)

// Result aggregates linting results from multiple paths
// Uses LinkedHashMap to preserve execution sequence and enable efficient result lookup
//
// Result 聚合来自多个路径的 linting 结果
// 使用 LinkedHashMap 保持执行顺序并实现高效的结果查找
type Result struct {
	resMap *linkedhashmap.Map[string, *golangcilint.Result] // Path to Result mapping with sequence preservation // 路径到结果的映射，保持顺序
}

// NewResult creates a Result instance with given result mapping
// Wraps the LinkedHashMap in a Result structure to enable method operations
//
// NewResult 创建 Result 实例，包含给定的结果映射
// 将 LinkedHashMap 包装在 Result 结构中以启用方法操作
func NewResult(resMap *linkedhashmap.Map[string, *golangcilint.Result]) *Result {
	return &Result{resMap: resMap}
}

// Success checks that each linting operation succeeded
// Returns false when paths have linting issues
//
// Success 检查每个 linting 操作是否成功
// 当路径有 linting 问题时返回 false
func (R *Result) Success() bool {
	for _, res := range R.resMap.Values() {
		if !res.Success() {
			return false
		}
	}
	return true
}

// GetMap returns the LinkedHashMap containing each result
// Provides access to path-result mappings in execution sequence
//
// GetMap 返回包含每个结果的 LinkedHashMap
// 按执行顺序提供路径-结果映射的访问
func (R *Result) GetMap() *linkedhashmap.Map[string, *golangcilint.Result] {
	return R.resMap
}

// DebugIssues displays diagnostic information about linting failures
// Shows both detailed and brief views of issues across each path
//
// DebugIssues 显示关于 linting 失败的诊断信息
// 显示每个路径问题的详细视图和简要视图
func (R *Result) DebugIssues() {
	if wrongCount := R.CountIssues(); wrongCount == 0 {
		eroticgo.GREEN.ShowMessage("SUCCESS")
		return
	}
	//首先显示详细的
	R.DebugIssues1()

	//接着显示缩略的
	R.DebugIssues2()
}

// CountIssues counts paths with linting failures
// Returns the count of failed paths
//
// CountIssues 计算有 linting 失败的路径数量
// 返回失败路径的数量
func (R *Result) CountIssues() int {
	var wrongCount int
	for _, res := range R.resMap.Values() {
		if !res.Success() {
			wrongCount++
		}
	}
	return wrongCount
}

// DebugIssues1 displays diagnostic information with complete issue details
// Shows execution commands, causes, and comprehensive issue listings
//
// DebugIssues1 显示包含完整问题详情的诊断信息
// 显示执行命令、原因和完整的问题列表
func (R *Result) DebugIssues1() {
	wrongCount := R.CountIssues()
	if wrongCount == 0 {
		return
	}
	eroticgo.AMBER.ShowMessage("FAILED", wrongCount, "WRONGS")
	{
		eroticgo.RED.ShowMessage("ERRORS:")

		var cnt int
		for idx, res := range R.resMap.Values() {
			fmt.Println(eroticgo.BLUE.Sprint("--"))
			fmt.Println(eroticgo.RED.Sprint("(", idx, ")", "path:", res.BasePath))
			fmt.Println(eroticgo.BLUE.Sprint("cd", res.BasePath, "&&", strings.Join([]string{"golangci-lint run --output.json.path=stdout --show-stats=false --timeout=5m0s"}, " ")))
			fmt.Println(eroticgo.BLUE.Sprint("--"))
			if res.Cause != nil {
				fmt.Println(eroticgo.RED.Sprint("command-execute-exception-cause:", res.Cause))
				cnt++
			} else if issues := res.Result.Issues; len(issues) > 0 {
				fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-issues:", len(issues)))
				cnt += len(issues)
				res.DebugIssues()
			} else {
				fmt.Println(eroticgo.GREEN.Sprint("success"))
			}
			fmt.Println(eroticgo.BLUE.Sprint("--"))
		}
		eroticgo.LIME.ShowMessage("FAILED", cnt, "ERRORS")
	}
	eroticgo.RED.ShowMessage("FAILED", wrongCount, "WRONGS")
}

// DebugIssues2 displays concise diagnostic information with issue summaries
// Shows abbreviated view with issue counts and file positions
//
// DebugIssues2 显示包含问题摘要的简洁诊断信息
// 显示缩略视图，包含问题计数和文件位置
func (R *Result) DebugIssues2() {
	wrongCount := R.CountIssues()
	if wrongCount == 0 {
		return
	}
	eroticgo.AMBER.ShowMessage("FAILED", wrongCount, "WRONGS")
	{
		eroticgo.RED.ShowMessage("ERRORS:")

		var cnt int
		for _, res := range R.resMap.Values() {
			if res.Cause != nil {
				fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", res.BasePath))
				fmt.Println(eroticgo.RED.Sprint("command-execute-exception-cause:", res.Cause))
				cnt++
			} else if issues := res.Result.Issues; len(issues) > 0 {
				fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", res.BasePath))
				fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-issues:", len(issues)))
				fmt.Println(eroticgo.RED.Sprint("--"))
				for _, issue := range issues {
					cnt++
					fmt.Println(eroticgo.RED.Sprint("pos:", filepath.Join(res.BasePath, issue.Pos.Filename)+":"+strconv.Itoa(issue.Pos.Line)+":"))
				}
			} else {
				continue
			}
			fmt.Println(eroticgo.RED.Sprint("--"))
		}
		eroticgo.LIME.ShowMessage("FAILED", cnt, "ERRORS")
	}
	eroticgo.RED.ShowMessage("FAILED", wrongCount, "WRONGS")
}
