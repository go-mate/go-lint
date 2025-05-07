package goworkcilint_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/goworkcilint"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestRun(t *testing.T) {
	roots := []string{
		runpath.PARENT.Path(),
	}
	results := goworkcilint.Run(osexec.NewExecConfig().WithDebug(), roots, 5*time.Minute)
	t.Log(neatjsons.S(results))

	goworkcilint.DebugIssues(roots, results)
}
