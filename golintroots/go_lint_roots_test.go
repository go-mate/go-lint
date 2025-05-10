package golintroots_test

import (
	"testing"
	"time"

	"github.com/go-mate/go-lint/golintroots"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestRun(t *testing.T) {
	roots := []string{
		runpath.PARENT.Path(),
	}
	resMap := golintroots.Run(osexec.NewExecConfig().WithDebug(), roots, 5*time.Minute)
	t.Log(neatjsons.S(resMap))

	golintroots.DebugIssues(roots, resMap)
}
