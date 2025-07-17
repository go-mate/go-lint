package gofrsflock_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/internal/gofrsflock"
	"github.com/stretchr/testify/require"
)

func TestLock(t *testing.T) {
	flock, err := gofrsflock.Lock(gofrsflock.CurrentPath(), time.Minute)
	require.NoError(t, err)
	time.Sleep(time.Second)
	require.NoError(t, flock.Unlock())
}

func TestRLock(t *testing.T) {
	flock, err := gofrsflock.RLock(gofrsflock.CurrentPath(), time.Minute)
	require.NoError(t, err)
	time.Sleep(time.Second)
	require.NoError(t, flock.Unlock())
}
