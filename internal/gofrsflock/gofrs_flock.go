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

func CurrentPath() string {
	return runpath.CurrentPath()
}

// Lock 尝试加锁，若锁文件不存在立即返回错误
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

// Unlock 解锁并释放文件对象
func Unlock(f *flock.Flock) error {
	return f.Unlock()
}
