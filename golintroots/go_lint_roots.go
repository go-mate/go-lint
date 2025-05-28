package golintroots

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/go-mate/go-lint/golangcilint"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/osexec"
)

type Result struct {
	resMap *linkedhashmap.Map[string, *golangcilint.Result]
}

func (R *Result) Success() bool {
	for _, res := range R.resMap.Values() {
		if !res.Success() {
			return false
		}
	}
	return true
}

func (R *Result) GetMap() *linkedhashmap.Map[string, *golangcilint.Result] {
	return R.resMap
}

func Run(execConfig *osexec.ExecConfig, roots []string, timeout time.Duration) *Result {
	var resMap = linkedhashmap.New[string, *golangcilint.Result]()
	for _, path := range roots {
		result := golangcilint.Run(execConfig, path, timeout)
		result.DebugIssues()
		resMap.Put(path, result)
	}
	return &Result{resMap: resMap}
}

func (R *Result) DebugIssues() {
	if wrongCount := CountIssues(R.resMap); wrongCount == 0 {
		eroticgo.GREEN.ShowMessage("SUCCESS")
		return
	}
	//首先显示详细的
	DebugIssues1(R.resMap)

	//接着显示缩略的
	DebugIssues2(R.resMap)
}

func CountIssues(resMap *linkedhashmap.Map[string, *golangcilint.Result]) int {
	var wrongCount int
	for _, res := range resMap.Values() {
		if !res.Success() {
			wrongCount++
		}
	}
	return wrongCount
}

func DebugIssues1(resMap *linkedhashmap.Map[string, *golangcilint.Result]) {
	wrongCount := CountIssues(resMap)
	if wrongCount == 0 {
		return
	}
	eroticgo.AMBER.ShowMessage("FAILED", wrongCount, "WRONGS")
	{
		eroticgo.RED.ShowMessage("ERRORS:")

		var cnt int
		for idx, res := range resMap.Values() {
			fmt.Println(eroticgo.BLUE.Sprint("--"))
			fmt.Println(eroticgo.RED.Sprint("(", idx, ")", "path:", res.BasePath))
			fmt.Println(eroticgo.BLUE.Sprint("cd", res.BasePath, "&&", strings.Join([]string{"golangci-lint run --output.json.path=stdout --show-stats=false --timeout=5m0s"}, " ")))
			fmt.Println(eroticgo.BLUE.Sprint("--"))
			if res.Cause != nil {
				fmt.Println(eroticgo.RED.Sprint("command-execute-exception-cause:", res.Cause))
				cnt++
			} else if issues := res.Result.Issues; len(issues) > 0 {
				fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-issues:", len(issues)))
				cnt += len(issues)
				res.DebugIssues()
			} else {
				fmt.Println(eroticgo.GREEN.Sprint("success"))
			}
			fmt.Println(eroticgo.BLUE.Sprint("--"))
		}
		eroticgo.LIME.ShowMessage("FAILED", cnt, "ERRORS")
	}
	eroticgo.RED.ShowMessage("FAILED", wrongCount, "WRONGS")
}

func DebugIssues2(resMap *linkedhashmap.Map[string, *golangcilint.Result]) {
	wrongCount := CountIssues(resMap)
	if wrongCount == 0 {
		return
	}
	eroticgo.AMBER.ShowMessage("FAILED", wrongCount, "WRONGS")
	{
		eroticgo.RED.ShowMessage("ERRORS:")

		var cnt int
		for _, res := range resMap.Values() {
			if res.Cause != nil {
				fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", res.BasePath))
				fmt.Println(eroticgo.RED.Sprint("command-execute-exception-cause:", res.Cause))
				cnt++
			} else if issues := res.Result.Issues; len(issues) > 0 {
				fmt.Println(eroticgo.RED.Sprint("(", cnt, ")", "path:", res.BasePath))
				fmt.Println(eroticgo.RED.Sprint("command-execute-wrong-issues:", len(issues)))
				fmt.Println(eroticgo.RED.Sprint("--"))
				for _, issue := range issues {
					cnt++
					fmt.Println(eroticgo.RED.Sprint("pos:", filepath.Join(res.BasePath, issue.Pos.Filename)+":"+strconv.Itoa(issue.Pos.Line)+":"))
				}
			} else {
				continue
			}
			fmt.Println(eroticgo.RED.Sprint("--"))
		}
		eroticgo.LIME.ShowMessage("FAILED", cnt, "ERRORS")
	}
	eroticgo.RED.ShowMessage("FAILED", wrongCount, "WRONGS")
}
