package golintroots

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

func (R *Result) Success() bool {
	return R.Reason == "" && R.Result.Success()
}

func Run(execConfig *osexec.ExecConfig, roots []string, timeout time.Duration) map[string]*Result {
	var resMap = map[string]*Result{}
	for _, path := range roots {
		result := golangcilint.Run(execConfig, path, timeout)
		if result.Reason != "" {
			fmt.Println(eroticgo.RED.Sprint("path:", path, "reason:", result.Reason))
			resMap[path] = &Result{
				Path:   path,
				Reason: result.Reason,
				Result: nil,
			}
			continue
		} else if len(result.Result.Issues) > 0 {
			fmt.Println(eroticgo.RED.Sprint("path:", path, "issues:"))
			golangcilint.DebugIssues(path, result)
			resMap[path] = &Result{
				Path:   path,
				Reason: "",
				Result: result,
			}
			continue
		} else {
			fmt.Println(eroticgo.GREEN.Sprint("path:", path, "return:", "success"))
			resMap[path] = &Result{
				Path:   path,
				Reason: "",
				Result: result,
			}
			continue
		}
	}
	return resMap
}

func DebugIssues(roots []string, resMap map[string]*Result) {
	var wrongCount int
	for _, result := range resMap {
		if !result.Success() {
			wrongCount++
		}
	}
	if wrongCount <= 0 {
		eroticgo.GREEN.ShowMessage("SUCCESS")
		return
	} else {
		//首先显示详细的
		eroticgo.PINK.ShowMessage("FAILED", wrongCount, "WRONGS")
		{
			eroticgo.RED.ShowMessage("ERRORS:")

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
					fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-reason:", res.Reason))
					cnt++
				} else if issues := res.Result.Result.Issues; len(issues) > 0 {
					fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-issues:", len(issues)))
					cnt += len(issues)
					golangcilint.DebugIssues(path, res.Result)
				} else {
					fmt.Println(eroticgo.GREEN.Sprint("success"))
				}
				fmt.Println(eroticgo.BLUE.Sprint("--"))
			}

			eroticgo.RED.ShowMessage("FAILED", cnt, "ERRORS")
		}

		//接着显示缩略的
		eroticgo.PINK.ShowMessage("FAILED", wrongCount, "WRONGS")
		{
			eroticgo.RED.ShowMessage("ERRORS:")

			var cnt int
			for _, path := range roots {
				res, ok := resMap[path]
				if !ok {
					continue
				}
				if res.Reason != "" {
					fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", path))
					fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-reason:", res.Reason))
					cnt++
				} else if issues := res.Result.Result.Issues; len(issues) > 0 {
					fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", path))
					fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-issues:", len(issues)))
					fmt.Println(eroticgo.RED.Sprint("--"))
					for _, issue := range issues {
						cnt++
						fmt.Println(eroticgo.RED.Sprint("pos:", filepath.Join(path, issue.Pos.Filename)+":"+strconv.Itoa(issue.Pos.Line)+":"))
					}
				} else {
					continue
				}
				fmt.Println(eroticgo.RED.Sprint("--"))
			}

			eroticgo.RED.ShowMessage("FAILED", cnt, "ERRORS")
		}
	}
}
