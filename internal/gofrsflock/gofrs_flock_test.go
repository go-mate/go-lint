// Package gofrsflock_test: Tests file locking with timeout protection
// Validates lock acquisition, holding, and release operations
// Ensures correct resource management with defer patterns
//
// gofrsflock_test: 测试带超时保护的文件锁定
// 验证锁定获取、持有和释放操作
// 确保使用 defer 模式进行正确的资源管理
package gofrsflock_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/internal/gofrsflock"
	"github.com/stretchr/testify/require"
)

// TestCurrentPath validates that CurrentPath returns consistent path
// This test verifies the lock file location mechanism
//
// TestCurrentPath 验证 CurrentPath 返回一致的路径
// 此测试验证锁文件位置机制
func TestCurrentPath(t *testing.T) {
	path := gofrsflock.CurrentPath()
	t.Log(path)
}

// TestLock validates read lock acquisition and release
// Tests basic file locking with timeout and safe cleanup
//
// TestLock 验证读锁定获取和释放
// 测试基础文件锁定和超时以及安全清理
func TestLock(t *testing.T) {
	flock, err := gofrsflock.Lock(gofrsflock.CurrentPath(), time.Minute)
	require.NoError(t, err)
	time.Sleep(time.Second)
	require.NoError(t, flock.Unlock())
}

// TestRLock validates explicit read lock acquisition and release
// Tests RLock function which has semantic naming as read lock operation
//
// TestRLock 验证显式读锁定获取和释放
// 测试 RLock 函数，这是语义化命名的读锁定操作
func TestRLock(t *testing.T) {
	flock, err := gofrsflock.RLock(gofrsflock.CurrentPath(), time.Minute)
	require.NoError(t, err)
	time.Sleep(time.Second)
	require.NoError(t, flock.Unlock())
}
