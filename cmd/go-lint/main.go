package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/emirpasic/gods/v2/sets/linkedhashset"
	"github.com/go-mate/go-lint/golint"
	"github.com/go-mate/go-lint/internal/utils"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	workRoot := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(workRoot))

	rootCmd := cobra.Command{
		Use:   "go-lint",
		Short: "go-lint",
		Long:  "go-lint",
		Run: func(cmd *cobra.Command, args []string) {
			runLint(workRoot)
		},
	}
	rootCmd.AddCommand(newLintRunCmd(workRoot))
	must.Done(rootCmd.Execute())
}

func newLintRunCmd(workRoot string) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "run",
		Long:  "run",
		Run: func(cmd *cobra.Command, args []string) {
			runLint(workRoot)
		},
	}
}

func runLint(workRoot string) {
	rootsHashSet := linkedhashset.New[string]()

	projectPath, shortMiddle, isGoModule := utils.GetProjectPath(workRoot)
	if !isGoModule {
		must.None(projectPath)
		must.None(shortMiddle)
	} else {
		rootsHashSet.Add(workRoot) //把当前目录添加到需要lint的目录里
	}

	//这里很有可能，当前目录下就是 go.mod 文件，就是把当前目录设置两次，因此使用 hash-set 去重复
	must.Done(filepath.Walk(workRoot, func(path string, info fs.FileInfo, err error) error {
		if retValue, isHide := isHidePath(info); isHide {
			return retValue
		}
		if !info.IsDir() && info.Name() == "go.mod" {
			if subRoot := filepath.Dir(path); osmustexist.IsRoot(subRoot) {
				rootsHashSet.Add(subRoot)
			}
			return nil
		}
		return nil
	}))
	zaplog.SUG.Debugln(neatjsons.S(rootsHashSet))

	//但是有些项目里是没有go文件的，比如空项目，或者大项目里只有子项目，而没有逻辑，因此需要去除
	rootsHashSet = rootsHashSet.Select(func(index int, value string) bool {
		zaplog.SUG.Debugln(index, value)
		return hasGoFiles(value)
	})
	zaplog.SUG.Debugln(neatjsons.S(rootsHashSet))

	result := golint.RootsRun(rootsHashSet.Values(), time.Minute*5)
	result.DebugIssues()
}

func hasGoFiles(root string) bool {
	exist := false
	must.Done(filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if retValue, isHide := isHidePath(info); isHide {
			return retValue
		}

		if info.IsDir() {
			//当遇到其它项目的 go.mod 时结束（由于传进来的就是项目根目录，因此要排除当前目录）
			if path != root && osmustexist.IsFile(filepath.Join(path, "go.mod")) {
				return filepath.SkipDir
			}
		} else {
			if filepath.Ext(info.Name()) == ".go" {
				exist = true
				return filepath.SkipAll
			}
		}
		return nil
	}))
	return exist
}

func isHidePath(info fs.FileInfo) (error, bool) {
	if info.IsDir() {
		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir, true
		}
	} else {
		if strings.HasPrefix(info.Name(), ".") {
			return nil, true
		}
	}
	return nil, false
}
