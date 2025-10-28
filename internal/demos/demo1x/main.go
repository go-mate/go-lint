package main

import (
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/yyle88/must"
	"github.com/yyle88/must/mustslice"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// go run main.go
//
// Demo showing low-level golangcilint.Run usage with direct result validation
// 演示底层 golangcilint.Run 的使用，带有直接的结果验证
func main() {
	projectPath := runpath.PARENT.Path()
	zaplog.SUG.Debugln(projectPath)

	// Execute golangci-lint with debug config and validate results
	// 使用调试配置执行 golangci-lint 并验证结果
	result := golangcilint.Run(osexec.NewCommandConfig().WithDebug(), projectPath, time.Minute*5)
	must.Done(result.Cause)
	mustslice.Have(result.Output)
	mustslice.None(result.Warnings)
	mustslice.None(result.Result.Issues)
}
