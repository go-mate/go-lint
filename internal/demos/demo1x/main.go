package main

import (
	"path/filepath"
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/yyle88/must/mustslice"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

// go run main.go
func main() {
	path := runpath.PARENT.Up(3)
	osmustexist.MustFile(filepath.Join(path, "go.mod"))

	res := rese.P1(golangcilint.Run(osexec.NewCommandConfig().WithDebugMode(true), path, time.Minute*5))
	mustslice.Have(res.Output)
	mustslice.None(res.Warnings)
	mustslice.None(res.Result.Issues)
}
