package main

import (
	"os"
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/go-mate/go-work/modulepath"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

// go-lint run
// go-lint run --debug
// go-lint run --debug=true
// go-lint run --debug=1
// go-lint run --debug=0
// go-lint
func main() {
	workPath := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(workPath))

	rootCmd := cobra.Command{
		Use:   "go-lint",
		Short: "go-lint",
		Long:  "go-lint",
		Run: func(cmd *cobra.Command, args []string) {
			run(workPath, false)
		},
	}
	rootCmd.AddCommand(newRunCmd(workPath))
	must.Done(rootCmd.Execute())
}

func newRunCmd(workPath string) *cobra.Command {
	var debugMode bool
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run",
		Long:  "run",
		Run: func(cmd *cobra.Command, args []string) {
			run(workPath, debugMode)
		},
	}
	cmd.Flags().BoolVarP(&debugMode, "debug", "", false, "enable debug mode")
	return cmd
}

func run(workPath string, debugMode bool) {
	golint.SetDebugMode(debugMode)
	modulePaths := modulepath.GetModulePaths(workPath, modulepath.NewOptions().WithDebugMode(debugMode))
	result := golint.RootsRun(modulePaths, time.Minute*5)
	result.DebugIssues()
}
