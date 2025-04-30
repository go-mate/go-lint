package goworkcilint

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/osexec"
)

type Result struct {
	Path   string
	Reason string
	Result *golangcilint.Result
}

func Run(execConfig *osexec.ExecConfig, roots []string, timeout time.Duration) map[string]*Result {
	var resMap = map[string]*Result{}
	for _, path := range roots {
		lintResult, err := golangcilint.Run(execConfig, path, timeout)
		if err != nil {
			fmt.Println(eroticgo.RED.Sprint("path:", path, "reason:", err))
			resMap[path] = &Result{
				Path:   path,
				Reason: err.Error(),
				Result: nil,
			}
			continue
		}
		if len(lintResult.Result.Issues) > 0 {
			fmt.Println(eroticgo.RED.Sprint("path:", path, "issues:"))
			golangcilint.DebugIssues(path, lintResult.Result.Issues)
			resMap[path] = &Result{
				Path:   path,
				Reason: "",
				Result: lintResult,
			}
			continue
		}
	}
	return resMap
}

func DebugIssues(roots []string, resMap map[string]*Result) {
	if len(resMap) > 0 {
		var wrongCount int

		{
			eroticgo.RED.ShowMessage("FAILED", len(resMap), "WRONGS")

			var cnt int
			var idx int
			for _, path := range roots {
				res, ok := resMap[path]
				if !ok {
					continue
				}
				idx++

				fmt.Println(eroticgo.BLUE.Sprint("--"))
				fmt.Println(eroticgo.RED.Sprint("(", idx, ")", "path:", path))
				fmt.Println(eroticgo.BLUE.Sprint("cd", path, "&&", strings.Join([]string{"golangci-lint run --output.json.path=stdout --show-stats=false --timeout=5m0s"}, " ")))
				fmt.Println(eroticgo.BLUE.Sprint("--"))
				if res.Reason != "" {
					cnt++
					fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-reason:", res.Reason))
				} else if issues := res.Result.Result.Issues; len(issues) > 0 {
					cnt += len(issues)
					golangcilint.DebugIssues(path, issues)
				}
				fmt.Println(eroticgo.BLUE.Sprint("--"))
			}

			eroticgo.RED.ShowMessage("FAILED", cnt, "ERRORS")

			wrongCount = cnt
		}

		if wrongCount > 0 {
			eroticgo.RED.ShowMessage("ERRORS:")

			var cnt int
			for _, path := range roots {
				res, ok := resMap[path]
				if !ok {
					continue
				}
				if res.Reason != "" {
					fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", path))
					cnt++
					fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-reason:", res.Reason))
				} else if issues := res.Result.Result.Issues; len(issues) > 0 {
					for _, issue := range issues {
						fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", path))
						cnt++
						fmt.Println(eroticgo.RED.Sprint("pos:", filepath.Join(path, issue.Pos.Filename)+":"+strconv.Itoa(issue.Pos.Line)+":"))
					}
				}
			}
		}
	} else {
		eroticgo.GREEN.ShowMessage("SUCCESS")
	}
}
