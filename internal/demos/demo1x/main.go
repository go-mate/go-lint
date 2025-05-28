package main

import (
	"path/filepath"
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/yyle88/must"
	"github.com/yyle88/must/mustslice"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// go run main.go
func main() {
	projectPath := runpath.PARENT.Up(3)
	osmustexist.MustFile(filepath.Join(projectPath, "go.mod"))
	zaplog.SUG.Debugln(projectPath)

	result := golangcilint.Run(osexec.NewCommandConfig().WithDebugMode(true), projectPath, time.Minute*5)
	must.Done(result.Cause)
	mustslice.Have(result.Output)
	mustslice.None(result.Warnings)
	mustslice.None(result.Result.Issues)
}
