// Package golangcilint_test: Tests core golangci-lint execution with JSON parsing
// Validates command execution, result parsing, and issue detection
// Uses file locks to prevent concurrent test conflicts
//
// golangcilint_test: 测试核心 golangci-lint 执行和 JSON 解析
// 验证命令执行、结果解析和问题检测
// 使用文件锁定防止并发测试冲突
package golangcilint_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/go-mate/go-lint/internal/gofrsflock"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// TestMain sets up the test environment
// Skips all tests when golangci-lint is not found in PATH
//
// TestMain 设置测试环境
// 当 PATH 中未找到 golangci-lint 时跳过所有测试
func TestMain(m *testing.M) {
	// Check whether golangci-lint is available
	// 检查 golangci-lint 是否可用
	path, err := exec.LookPath("golangci-lint")
	if err != nil {
		zaplog.SUG.Warnln("golangci-lint not found in PATH, skipping tests")
		os.Exit(0)
	}
	zaplog.LOG.Debug("golangci-lint found in PATH", zap.String("path", path))

	m.Run()
}

// TestRun validates golangci-lint execution with debug mode
// Ensures execution completes, result parsing works, and issue reporting functions
//
// TestRun 验证 golangci-lint 执行和调试模式
// 确保执行完成、结果解析工作和问题报告功能
func TestRun(t *testing.T) {
	flock := rese.P1(gofrsflock.Lock(gofrsflock.CurrentPath(), 5*time.Minute+30*time.Second))
	defer rese.F0(flock.Unlock)

	root := runpath.PARENT.Path()
	result := golangcilint.Run(osexec.NewExecConfig().WithDebug(), root, 5*time.Minute)
	require.NoError(t, result.Cause)
	t.Log(neatjsons.S(result))

	result.DebugIssues()
	require.True(t, result.Success())
}
