// Package gofrsflock: File locking utilities using gofrs/flock package
// Provides simplified interface for read locks with timeout protection
// Used in test scenarios to prevent concurrent golangci-lint execution conflicts
//
// gofrsflock: 使用 gofrs/flock 包的文件锁定实用程序
// 提供带超时保护的读锁定的简化接口
// 在测试场景中使用，防止并发 golangci-lint 执行冲突
package gofrsflock

import (
	"context"
	"os"
	"time"

	"github.com/gofrs/flock"
	"github.com/yyle88/erero"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
)

// CurrentPath returns the path of the current executing package
// Convenience wrapper around runpath.CurrentPath
//
// CurrentPath 返回当前执行包的路径
// runpath.CurrentPath 的便捷包装器
func CurrentPath() string {
	return runpath.CurrentPath()
}

// Lock attempts to acquire a read lock on the specified file with timeout
// Returns error if lock file does not exist or cannot be acquired within timeout
//
// Lock 尝试在指定文件上获取读锁定，带有超时
// 如果锁文件不存在或无法在超时内获取则返回错误
func Lock(lockPath string, timeout time.Duration) (*flock.Flock, error) {
	f := flock.New(osmustexist.FILE(lockPath))
	locked, err := f.TryRLockContext(context.Background(), timeout)
	if err != nil {
		return nil, erero.Wro(err)
	}
	if !locked {
		return nil, os.ErrDeadlineExceeded
	}
	return f, nil
}

// RLock attempts to acquire a read lock on the specified file with timeout
// Same as Lock but explicitly named to indicate read locking behavior
//
// RLock 尝试在指定文件上获取读锁定，带有超时
// 与 Lock 相同，但明确命名以指示读锁定行为
func RLock(lockPath string, timeout time.Duration) (*flock.Flock, error) {
	f := flock.New(osmustexist.FILE(lockPath))
	locked, err := f.TryRLockContext(context.Background(), timeout)
	if err != nil {
		return nil, erero.Wro(err)
	}
	if !locked {
		return nil, os.ErrDeadlineExceeded
	}
	return f, nil
}

// Unlock releases the file lock and frees associated resources
// Should be called when the protected operation completes
//
// Unlock 释放文件锁定并释放相关资源
// 应在受保护操作完成时调用
func Unlock(f *flock.Flock) error {
	return f.Unlock()
}
