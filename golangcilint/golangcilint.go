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
	LintResult *printers.JSONResult
	Warnings   []string
	RawMessage []byte
}

func Run(execConfig *osexec.ExecConfig, path string, timeout time.Duration) (*Result, error) {
	rawMessage, err := execConfig.ShallowClone().WithPath(path).Exec("golangci-lint", "run", "--output.json.path=stdout", "--show-stats=false", "--timeout="+timeout.String())
	if err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))

		//假如能够顺利的转化为 json 结果，返回 issues
		if res := parseMessage(rawMessage); res != nil {
			return debugMessage(res), nil
		}

		//假如排除 warning 以后能够转化为 json 结果，也返回 issues
		if res := parseIgnoreWarningMessage(strings.Split(string(rawMessage), "\n")); res != nil {
			return debugMessage(res), nil
		}

		//当代码有大错时，没法返回细致的 issues，这时就需要把这个大错显示出来
		zaplog.SUG.Errorln("message:", string(rawMessage))
		return nil, erero.Wro(err)
	}
	resultLint := &printers.JSONResult{}
	must.Done(json.Unmarshal(rawMessage, resultLint))
	return debugMessage(&Result{
		LintResult: resultLint,
		Warnings:   nil,
		RawMessage: rawMessage,
	}), nil
}

func parseMessage(rawMessage []byte) *Result {
	resultLint := &printers.JSONResult{}
	if err := json.Unmarshal(rawMessage, resultLint); err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))
		return nil
	}
	return &Result{
		LintResult: resultLint,
		Warnings:   nil,
		RawMessage: rawMessage,
	}
}

func debugMessage(res *Result) *Result {
	zaplog.SUG.Debugln("message:", neatjsons.SxB(res.RawMessage))
	for _, warning := range res.Warnings {
		zaplog.SUG.Warnln("warning:", warning)
	}
	for _, issue := range res.LintResult.Issues {
		zaplog.SUG.Errorln(issue.Text)
	}
	return res
}

func parseIgnoreWarningMessage(msgLines []string) *Result {
	var warnings []string
	var jsonLineCount = 0
	var jsonLine string
	var unknownOutput = false
	for _, msg := range msgLines {
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
	resultLint := &printers.JSONResult{}
	rawMessage := []byte(jsonLine)
	if err := json.Unmarshal(rawMessage, resultLint); err != nil {
		zaplog.LOG.Debug("reason:", zap.Error(err))
		return nil
	}
	return &Result{
		LintResult: resultLint,
		Warnings:   warnings,
		RawMessage: rawMessage,
	}
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
