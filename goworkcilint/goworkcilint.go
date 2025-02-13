package goworkcilint

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/osexec"
)

type Result struct {
	Path       string
	CmdExecErr error
	RawMessage []byte
	Warnings   []string
	LintIssues []result.Issue
}

func Run(execConfig *osexec.ExecConfig, roots []string, timeout time.Duration) map[string]*Result {
	var resMap = map[string]*Result{}
	for _, projectPath := range roots {
		res, err := golangcilint.Run(execConfig, projectPath, timeout)
		if err != nil {
			fmt.Println(eroticgo.RED.Sprint("path:", projectPath, "reason:", err))
			resMap[projectPath] = &Result{
				Path:       projectPath,
				CmdExecErr: err,
				RawMessage: []byte{},
				Warnings:   []string{},
				LintIssues: nil,
			}
			continue
		}
		if len(res.LintResult.Issues) > 0 {
			fmt.Println(eroticgo.RED.Sprint("path:", projectPath, "issues:"))
			golangcilint.DebugIssues(projectPath, res.LintResult.Issues)
			resMap[projectPath] = &Result{
				Path:       projectPath,
				CmdExecErr: nil,
				RawMessage: res.RawMessage,
				Warnings:   res.Warnings,
				LintIssues: res.LintResult.Issues,
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
				fmt.Println(eroticgo.BLUE.Sprint("cd", path, "&&", strings.Join([]string{"golangci-lint", "run", "--out-format", "json"}, " ")))
				fmt.Println(eroticgo.BLUE.Sprint("--"))
				if res.CmdExecErr != nil {
					cnt++
					fmt.Println(eroticgo.RED.Sprint("command-execute-error-message:", res.CmdExecErr))
				} else if issues := res.LintIssues; len(issues) > 0 {
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
				if res.CmdExecErr != nil {
					fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", path))
					cnt++
					fmt.Println(eroticgo.RED.Sprint("command-execute-error-message:", res.CmdExecErr))
				} else if issues := res.LintIssues; len(issues) > 0 {
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
