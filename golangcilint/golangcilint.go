// Package golangcilint: Core golangci-lint execution engine with advanced result processing
// Provides direct integration with golangci-lint CLI tool with JSON result parsing
// Handles complex error scenarios, warning processing, and detailed issue reporting
//
// golangcilint: 核心 golangci-lint 执行引擎，支持高级结果处理
// 提供与 golangci-lint CLI 工具的直接集成，支持 JSON 结果解析
// 处理复杂错误场景、警告处理和详细问题报告
package golangcilint

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// Result represents the comprehensive output of a golangci-lint execution
// Contains execution status, parsed issues, warnings, and raw output for analysis
//
// Result 表示 golangci-lint 执行的综合输出
// 包含执行状态、解析的问题、警告和用于分析的原始输出
type Result struct {
	BasePath string               // Base path where linting was executed // 执行 linting 的基础路径
	Cause    error                // Cause of the command "golangci-lint run" failure // 命令 "golangci-lint run" 失败的原因
	Result   *printers.JSONResult // Parsed JSON result from golangci-lint // 来自 golangci-lint 的解析 JSON 结果
	Warnings []string             // Warning messages during execution // 执行期间的警告消息
	Output   json.RawMessage      // Raw JSON output for debugging // 用于调试的原始 JSON 输出
}

// Success returns true if linting completed without errors, warnings, or issues
// Indicates clean code that passes all golangci-lint checks
//
// Success 在 linting 完成且无错误、警告或问题时返回 true
// 表示通过所有 golangci-lint 检查的干净代码
func (R *Result) Success() bool {
	return R.Cause == nil && len(R.Warnings) == 0 && len(R.Result.Issues) == 0
}

// Run executes golangci-lint with JSON output parsing and comprehensive error handling
// Handles multiple failure scenarios: complete failures, partial results with warnings, and clean success
// Returns structured results suitable for further processing and debugging
//
// Run 执行 golangci-lint，支持 JSON 输出解析和全面错误处理
// 处理多种失败场景：完全失败、带警告的部分结果和干净成功
// 返回适合进一步处理和调试的结构化结果
func Run(execConfig *osexec.ExecConfig, path string, timeout time.Duration) *Result {
	// Execute golangci-lint with JSON output and timeout protection
	// 执行 golangci-lint，使用 JSON 输出和超时保护
	rawMessage, err := execConfig.SubConfig(path).Exec("golangci-lint", "run", "--output.json.path=stdout", "--show-stats=false", "--timeout="+timeout.String())
	if err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))

		// Try to parse JSON result even when command fails (partial success scenario)
		// 即使命令失败也尝试解析 JSON 结果（部分成功场景）
		if res := parseMessage(rawMessage); res != nil {
			res.BasePath = path
			return debugMessage(execConfig, res)
		}

		// Try parsing after filtering out warning messages
		// 过滤警告消息后尝试解析
		if res := parseSkipWarningMessage(rawMessage); res != nil {
			res.BasePath = path
			return debugMessage(execConfig, res)
		}

		// Complete failure: cannot extract structured issues, show raw error
		// 完全失败：无法提取结构化问题，显示原始错误
		zaplog.SUG.Errorln("message:", string(rawMessage))
		return &Result{
			BasePath: path,
			Cause:    err,
			Result:   nil,
			Warnings: nil,
			Output:   nil,
		}
	}
	// Success path: parse clean JSON result
	// 成功路径：解析干净的 JSON 结果
	lintResult := &printers.JSONResult{}
	must.Done(json.Unmarshal(rawMessage, lintResult))
	res := &Result{
		BasePath: path,
		Cause:    nil,
		Result:   lintResult,
		Warnings: nil,
		Output:   rawMessage,
	}
	return debugMessage(execConfig, res)
}

// parseMessage attempts to parse raw output as clean JSON result
// Used for scenarios where golangci-lint fails but still produces valid JSON
//
// parseMessage 尝试将原始输出解析为干净的 JSON 结果
// 用于 golangci-lint 失败但仍产生有效 JSON 的场景
func parseMessage(rawMessage []byte) *Result {
	lintResult := &printers.JSONResult{}
	if err := json.Unmarshal(rawMessage, lintResult); err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))
		return nil
	}
	return &Result{
		Cause:    nil,
		Result:   lintResult,
		Warnings: nil,
		Output:   rawMessage,
	}
}

// debugMessage processes and logs detailed result information based on debug configuration
// Outputs formatted JSON, warnings, and issues for comprehensive analysis
//
// debugMessage 基于调试配置处理并记录详细结果信息
// 输出格式化的 JSON、警告和问题进行全面分析
func debugMessage(config *osexec.ExecConfig, res *Result) *Result {
	if config.IsShowOutputs() {
		zaplog.SUG.Debugln("message:", neatjsons.SxB(res.Output))
	}
	for _, warning := range res.Warnings {
		zaplog.SUG.Warnln("warning:", warning)
	}
	for _, issue := range res.Result.Issues {
		zaplog.SUG.Errorln(issue.Text)
	}
	return res
}

// parseSkipWarningMessage extracts JSON result from mixed output containing warnings
// Separates warning lines from JSON content for clean result parsing
//
// parseSkipWarningMessage 从包含警告的混合输出中提取 JSON 结果
// 将警告行与 JSON 内容分离以进行干净的结果解析
func parseSkipWarningMessage(rawMessage []byte) *Result {
	var warnings []string
	var jsonLineCount = 0
	var jsonLine string
	var unknownOutput = false
	for _, msg := range strings.Split(string(rawMessage), "\n") {
		if msg == "" {
			continue
		} else if strings.HasPrefix(msg, "level=warning") {
			zaplog.SUG.Warnln(msg)
			warnings = append(warnings, msg)
		} else if strings.HasPrefix(msg, "{") && strings.HasSuffix(msg, "}") {
			jsonLineCount++
			jsonLine = msg
		} else {
			unknownOutput = true
			zaplog.SUG.Errorln(msg)
		}
	}
	if unknownOutput {
		return nil
	}
	if jsonLineCount != 1 {
		return nil
	}
	lintResult := &printers.JSONResult{}
	jsonOutput := []byte(jsonLine)
	if err := json.Unmarshal(jsonOutput, lintResult); err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))
		return nil
	}
	return &Result{
		Cause:    nil,
		Result:   lintResult,
		Warnings: warnings,
		Output:   jsonOutput,
	}
}

// ShowCommandMessage displays the command that will be executed
// Shows the cd and golangci-lint command before execution
//
// ShowCommandMessage 显示将要执行的命令
// 在执行前显示 cd 和 golangci-lint 命令
func ShowCommandMessage(path string) {
	fmt.Println(eroticgo.BLUE.Sprint("--"))
	fmt.Println(eroticgo.BLUE.Sprint("cd", path, "&&", "golangci-lint run"))
	fmt.Println(eroticgo.BLUE.Sprint("--"))
}

// ShowOutlineMessage displays a brief result message immediately after execution
// Shows simple success/failure status without detailed issue breakdown
//
// ShowOutlineMessage 在执行后立即显示简要结果消息
// 显示简单的成功/失败状态，不含详细问题分解
func (R *Result) ShowOutlineMessage() {
	fmt.Println(eroticgo.BLUE.Sprint("--"))
	if R.Cause != nil {
		fmt.Println(eroticgo.PINK.Sprint("CAUSE:", R.Cause))
	} else if len(R.Result.Issues) > 0 {
		fmt.Println(eroticgo.PINK.Sprint("ISSUES:", len(R.Result.Issues)))
	} else {
		fmt.Println(eroticgo.AQUA.Sprint("SUCCESS"))
	}
	fmt.Println(eroticgo.BLUE.Sprint("--"))
}

// DebugIssues displays comprehensive, colorized output of linting results
// Shows command line, execution status, and detailed issue information with file positions
// Provides user-friendly visualization of linting problems for easy navigation
//
// DebugIssues 显示 linting 结果的综合彩色输出
// 显示命令行、执行状态和带文件位置的详细问题信息
// 为用户提供友好的 linting 问题可视化，便于导航
func (R *Result) DebugIssues() {
	commandLine := "cd " + must.Nice(R.BasePath) + " && golangci-lint run"

	fmt.Println(eroticgo.BLUE.Sprint("--"))
	if R.Cause != nil {
		// Show execution failure with cause
		// 显示带原因的执行失败
		fmt.Println(eroticgo.RED.Sprint(commandLine, "->", "exception-cause:", R.Cause))
	} else if issues := R.Result.Issues; len(issues) > 0 {
		// Show detailed issue breakdown for each linting problem
		// 为每个 linting 问题显示详细问题分解
		fmt.Println(eroticgo.RED.Sprint(commandLine), "->", "warning")
		for _, issueItem := range issues {
			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			// Absolute file path with line number
			// 绝对文件路径和行号
			fmt.Println(eroticgo.RED.Sprint("pos:", filepath.Join(must.Nice(R.BasePath), issueItem.Pos.Filename)+":"+strconv.Itoa(issueItem.Pos.Line)+":"))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			// Relative file path with line number
			// 相对文件路径和行号
			fmt.Println(eroticgo.RED.Sprint("pos:", issueItem.Pos.Filename+":"+strconv.Itoa(issueItem.Pos.Line)+":"))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			// Issue text message
			// 问题文本消息
			fmt.Println(eroticgo.RED.Sprint(strings.TrimSpace(issueItem.Text)))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			// Detailed issue description
			// 详细问题描述
			fmt.Println(eroticgo.RED.Sprint(strings.TrimSpace(issueItem.Description())))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			// Full JSON representation for debugging
			// 用于调试的完整 JSON 表示
			res := neatjsons.S(issueItem)
			fmt.Println(eroticgo.RED.Sprint("res:", res))
		}
	} else {
		// Clean success case
		// 干净成功情况
		fmt.Println(eroticgo.GREEN.Sprint(commandLine, "->", "success"))
	}
	fmt.Println(eroticgo.BLUE.Sprint("--"))
}
