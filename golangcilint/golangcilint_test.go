package golangcilint_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestRun(t *testing.T) {
	root := runpath.PARENT.Path()
	result := golangcilint.Run(osexec.NewExecConfig().WithDebug(), root, 5*time.Minute)
	require.NoError(t, result.Cause)
	t.Log(neatjsons.S(result))

	result.DebugIssues()
	require.True(t, result.Success())
}
