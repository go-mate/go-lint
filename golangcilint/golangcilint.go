package golangcilint

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

type Result struct {
	BasePath string
	Cause    error // cause of the command "golangci-lint run" failure
	Result   *printers.JSONResult
	Warnings []string
	Output   json.RawMessage
}

func (R *Result) Success() bool {
	return R.Cause == nil && len(R.Warnings) == 0 && len(R.Result.Issues) == 0
}

func Run(execConfig *osexec.ExecConfig, path string, timeout time.Duration) *Result {
	rawMessage, err := execConfig.SubConfig(path).Exec("golangci-lint", "run", "--output.json.path=stdout", "--show-stats=false", "--timeout="+timeout.String())
	if err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))

		//假如能够顺利的转化为 json 结果，返回 issues
		if res := parseMessage(rawMessage); res != nil {
			res.BasePath = path
			return debugMessage(res)
		}

		//假如排除 warning 以后能够转化为 json 结果，也返回 issues
		if res := parseSkipWarningMessage(rawMessage); res != nil {
			res.BasePath = path
			return debugMessage(res)
		}

		//当代码有大错时，没法返回细致的 issues，这时就需要把这个大错显示出来
		zaplog.SUG.Errorln("message:", string(rawMessage))
		return &Result{
			BasePath: path,
			Cause:    err,
			Result:   nil,
			Warnings: nil,
			Output:   nil,
		}
	}
	lintResult := &printers.JSONResult{}
	must.Done(json.Unmarshal(rawMessage, lintResult))
	return debugMessage(&Result{
		BasePath: path,
		Cause:    nil,
		Result:   lintResult,
		Warnings: nil,
		Output:   rawMessage,
	})
}

func parseMessage(rawMessage []byte) *Result {
	lintResult := &printers.JSONResult{}
	if err := json.Unmarshal(rawMessage, lintResult); err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))
		return nil
	}
	return &Result{
		Cause:    nil,
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
		Cause:    nil,
		Result:   lintResult,
		Warnings: warnings,
		Output:   jsonOutput,
	}
}

func (R *Result) DebugIssues() {
	commandLine := "cd " + must.Nice(R.BasePath) + " && golangci-lint run"

	fmt.Println(eroticgo.BLUE.Sprint("--"))
	if R.Cause != nil {
		fmt.Println(eroticgo.RED.Sprint(commandLine, "->", "exception-cause:", R.Cause))
	} else if issues := R.Result.Issues; len(issues) > 0 {
		fmt.Println(eroticgo.RED.Sprint(commandLine), "->", "warning")
		for _, issueItem := range issues {
			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			fmt.Println(eroticgo.RED.Sprint("pos:", filepath.Join(must.Nice(R.BasePath), issueItem.Pos.Filename)+":"+strconv.Itoa(issueItem.Pos.Line)+":"))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			fmt.Println(eroticgo.RED.Sprint("pos:", issueItem.Pos.Filename+":"+strconv.Itoa(issueItem.Pos.Line)+":"))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			fmt.Println(eroticgo.RED.Sprint(strings.TrimSpace(issueItem.Text)))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			fmt.Println(eroticgo.RED.Sprint(strings.TrimSpace(issueItem.Description())))

			fmt.Println(eroticgo.YELLOW.Sprint("--"))

			res := neatjsons.S(issueItem)

			fmt.Println(eroticgo.RED.Sprint("res:", res))
		}
	} else {
		fmt.Println(eroticgo.GREEN.Sprint("--"))
		fmt.Println(eroticgo.GREEN.Sprint("--"))
		fmt.Println(eroticgo.GREEN.Sprint(commandLine, "->", "success"))
		fmt.Println(eroticgo.GREEN.Sprint("--"))
		fmt.Println(eroticgo.GREEN.Sprint("--"))
	}
	fmt.Println(eroticgo.BLUE.Sprint("--"))
}
