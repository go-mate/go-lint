// Package golint_test: Tests the golint execution functions with file locking
// Validates single path, multiple paths, and batch execution modes
// Uses file locks to prevent concurrent test conflicts
//
// golint_test: 测试 golint 执行函数，带有文件锁定
// 验证单路径、多路径和批量执行模式
// 使用文件锁定防止并发测试冲突
package golint_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/go-mate/go-lint/internal/gofrsflock"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// projectPath holds the test project root path
// Initialized in TestMain and used across test cases
//
// projectPath 保存测试项目根路径
// 在 TestMain 中初始化并在测试用例中使用
var projectPath string

// TestMain sets up the test environment and initializes projectPath
// Validates that go.mod exists before running tests
//
// TestMain 设置测试环境并初始化 projectPath
// 在运行测试前验证 go.mod 存在
func TestMain(m *testing.M) {
	path := runpath.PARENT.Up(1)
	osmustexist.MustFile(filepath.Join(path, "go.mod"))
	zaplog.SUG.Debugln(path)

	projectPath = path
	m.Run()
}

// TestRun validates single path linting execution
// Uses file lock to prevent concurrent golangci-lint conflicts
//
// TestRun 验证单路径 linting 执行
// 使用文件锁定防止并发 golangci-lint 冲突
func TestRun(t *testing.T) {
	flock := rese.P1(gofrsflock.Lock(gofrsflock.CurrentPath(), 5*time.Minute+30*time.Second))
	defer rese.F0(flock.Unlock)

	result := golint.Run(projectPath, 5*time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
}

// TestWorksRun validates multi-project linting execution
// Tests WorksRun function with project path slice
//
// TestWorksRun 验证多项目 linting 执行
// 测试 WorksRun 函数使用项目路径切片
func TestWorksRun(t *testing.T) {
	flock := rese.P1(gofrsflock.Lock(gofrsflock.CurrentPath(), 5*time.Minute+30*time.Second))
	defer rese.F0(flock.Unlock)

	projectPaths := []string{projectPath}
	result := golint.WorksRun(projectPaths, 5*time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
}

// TestRootsRun validates multiple roots linting execution
// Tests RootsRun function with root path slice
//
// TestRootsRun 验证多根路径 linting 执行
// 测试 RootsRun 函数使用根路径切片
func TestRootsRun(t *testing.T) {
	flock := rese.P1(gofrsflock.Lock(gofrsflock.CurrentPath(), 5*time.Minute+30*time.Second))
	defer rese.F0(flock.Unlock)

	result := golint.RootsRun([]string{projectPath}, 5*time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
}

// TestBatchRun validates batch linting execution with custom configuration
// Tests BatchRun function with debug mode and validates result count
//
// TestBatchRun 验证带自定义配置的批量 linting 执行
// 测试 BatchRun 函数使用调试模式并验证结果计数
func TestBatchRun(t *testing.T) {
	flock := rese.P1(gofrsflock.Lock(gofrsflock.CurrentPath(), 5*time.Minute+30*time.Second))
	defer rese.F0(flock.Unlock)

	roots := []string{
		runpath.PARENT.Path(),
	}
	result := golint.BatchRun(osexec.NewExecConfig().WithDebug(), roots, 5*time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
	require.Equal(t, result.GetMap().Size(), len(roots))
}
