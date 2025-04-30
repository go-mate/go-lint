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
	Result   *printers.JSONResult
	Warnings []string
	Output   []byte
}

func Run(execConfig *osexec.ExecConfig, path string, timeout time.Duration) (*Result, error) {
	rawMessage, err := execConfig.GetSubClone(path).Exec("golangci-lint", "run", "--output.json.path=stdout", "--show-stats=false", "--timeout="+timeout.String())
	if err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))

		//假如能够顺利的转化为 json 结果，返回 issues
		if res := parseMessage(rawMessage); res != nil {
			return debugMessage(res), nil
		}

		//假如排除 warning 以后能够转化为 json 结果，也返回 issues
		if res := parseSkipWarningMessage(rawMessage); res != nil {
			return debugMessage(res), nil
		}

		//当代码有大错时，没法返回细致的 issues，这时就需要把这个大错显示出来
		zaplog.SUG.Errorln("message:", string(rawMessage))
		return nil, erero.Wro(err)
	}
	lintResult := &printers.JSONResult{}
	must.Done(json.Unmarshal(rawMessage, lintResult))
	return debugMessage(&Result{
		Result:   lintResult,
		Warnings: nil,
		Output:   rawMessage,
	}), nil
}

func parseMessage(rawMessage []byte) *Result {
	lintResult := &printers.JSONResult{}
	if err := json.Unmarshal(rawMessage, lintResult); err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))
		return nil
	}
	return &Result{
		Result:   lintResult,
		Warnings: nil,
		Output:   rawMessage,
	}
}

func debugMessage(res *Result) *Result {
	zaplog.SUG.Debugln("message:", neatjsons.SxB(res.Output))
	for _, warning := range res.Warnings {
		zaplog.SUG.Warnln("warning:", warning)
	}
	for _, issue := range res.Result.Issues {
		zaplog.SUG.Errorln(issue.Text)
	}
	return res
}

func parseSkipWarningMessage(rawMessage []byte) *Result {
	var warnings []string
	var jsonLineCount = 0
	var jsonLine string
	var unknownOutput = false
	for _, msg := range strings.Split(string(rawMessage), "\n") {
		if msg == "" {
			continue
		} else if strings.HasPrefix(msg, "level=warning") {
			zaplog.SUG.Warnln(msg)
			warnings = append(warnings, msg)
		} else if strings.HasPrefix(msg, "{") && strings.HasSuffix(msg, "}") {
			jsonLineCount++
			jsonLine = msg
		} else {
			unknownOutput = true
			zaplog.SUG.Errorln(msg)
		}
	}
	if unknownOutput {
		return nil
	}
	if jsonLineCount != 1 {
		return nil
	}
	lintResult := &printers.JSONResult{}
	jsonOutput := []byte(jsonLine)
	if err := json.Unmarshal(jsonOutput, lintResult); err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))
		return nil
	}
	return &Result{
		Result:   lintResult,
		Warnings: warnings,
		Output:   jsonOutput,
	}
}

func DebugIssues(root string, issues []result.Issue) {
	commandLine := "cd " + root + " && golangci-lint run"
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
