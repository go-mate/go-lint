package golangcilint

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

type Result struct {
	RawMessage []byte
	Warnings   []string
	LintResult *printers.JSONResult
}

func Run(execConfig *osexec.ExecConfig, path string, timeout time.Duration) (*Result, error) {
	lintResult := &printers.JSONResult{}
	outputData, err := execConfig.ShallowClone().WithPath(path).Exec("golangci-lint", "run", "--out-format", "json", "--timeout", timeout.String())
	if err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))

		//假如能够顺利的转化为结果，返回 issues
		if json.Unmarshal(outputData, lintResult) == nil {
			res := &Result{
				RawMessage: outputData,
				Warnings:   nil,
				LintResult: lintResult,
			}
			debugLintMessage(res)
			return res, nil
		}

		//假如排除 warning 以后能够转化为 json 结果，也返回 issues
		if res := parseLintOutput(strings.Split(string(outputData), "\n")); res != nil {
			debugLintMessage(res)
			return res, nil
		}

		//当代码有大错时，没法返回细致的 issues，这时就需要把这个大错显示出来
		zaplog.SUG.Errorln("output:", string(outputData))
		return nil, erero.Wro(err)
	}
	must.Done(json.Unmarshal(outputData, lintResult))
	res := &Result{
		RawMessage: outputData,
		Warnings:   nil,
		LintResult: lintResult,
	}
	debugLintMessage(res)
	return res, nil
}

func debugLintMessage(res *Result) {
	zaplog.SUG.Debugln("message:", neatjsons.SxB(res.RawMessage))
	for _, warning := range res.Warnings {
		zaplog.SUG.Warnln("warning:", warning)
	}
	for _, issue := range res.LintResult.Issues {
		zaplog.SUG.Errorln(issue.Text)
	}
}

func parseLintOutput(outputLines []string) *Result {
	var warnings []string
	var jsonLineCount = 0
	var jsonLine string
	var unknown = false
	for _, textLine := range outputLines {
		if textLine == "" {
			continue
		} else if strings.HasPrefix(textLine, "level=warning") {
			zaplog.SUG.Warnln(textLine)
			warnings = append(warnings, textLine)
		} else if strings.HasPrefix(textLine, "{") && strings.HasSuffix(textLine, "}") {
			jsonLineCount++
			jsonLine = textLine
		} else {
			unknown = true
			zaplog.SUG.Errorln(textLine)
		}
	}
	if !unknown && jsonLineCount == 1 {
		res := &printers.JSONResult{}
		rawOutput := []byte(jsonLine)
		if json.Unmarshal(rawOutput, res) == nil {
			return &Result{
				RawMessage: rawOutput,
				Warnings:   warnings,
				LintResult: res,
			}
		}
	}
	return nil
}

func DebugIssues(root string, issues []result.Issue) {
	commandLine := fmt.Sprintf("cd %s && golangci-lint run", root)
	if len(issues) > 0 {
		fmt.Println(eroticgo.BLUE.Sprint("--"))
		fmt.Println(eroticgo.RED.Sprint(commandLine), "->", "warning")
		for _, issue := range issues {
			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			fmt.Println(eroticgo.RED.Sprint("pos:", filepath.Join(root, issue.Pos.Filename)+":"+strconv.Itoa(issue.Pos.Line)+":"))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			fmt.Println(eroticgo.RED.Sprint("pos:", issue.Pos.Filename+":"+strconv.Itoa(issue.Pos.Line)+":"))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			fmt.Println(eroticgo.RED.Sprint(strings.TrimSpace(issue.Text)))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			fmt.Println(eroticgo.RED.Sprint(strings.TrimSpace(issue.Description())))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			res := neatjsons.S(issue)

			fmt.Println(eroticgo.RED.Sprint("res:", res))
		}
		fmt.Println(eroticgo.BLUE.Sprint("--"))
	} else {
		fmt.Println(eroticgo.BLUE.Sprint("--"))
		fmt.Println(eroticgo.GREEN.Sprint(commandLine, "->", "success"))
		fmt.Println(eroticgo.BLUE.Sprint("--"))
	}
}
