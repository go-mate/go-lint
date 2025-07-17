package golintroots_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/golintroots"
	"github.com/go-mate/go-lint/internal/gofrsflock"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

func TestRun(t *testing.T) {
	flock := rese.P1(gofrsflock.Lock(gofrsflock.CurrentPath(), 5*time.Minute+30*time.Second))
	defer rese.F0(flock.Unlock)

	roots := []string{
		runpath.PARENT.Path(),
	}
	result := golintroots.Run(osexec.NewExecConfig().WithDebug(), roots, 5*time.Minute)
	result.DebugIssues()
	require.True(t, result.Success())
	require.Equal(t, result.GetMap().Size(), len(roots))
}
