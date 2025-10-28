[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-mate/go-lint/release.yml?branch=main&label=BUILD)](https://github.com/go-mate/go-lint/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-mate/go-lint)](https://pkg.go.dev/github.com/go-mate/go-lint)
[![Coverage Status](https://img.shields.io/coveralls/github/go-mate/go-lint/main.svg)](https://coveralls.io/github/go-mate/go-lint?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-mate/go-lint.svg)](https://github.com/go-mate/go-lint/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mate/go-lint)](https://goreportcard.com/report/github.com/go-mate/go-lint)

<p align="center">
  <img
    alt="golangci-lint logo"
    src="assets/golangci-lint-logo.jpeg"
    style="max-height: 500px; width: auto; max-width: 100%;"
  />
</p>
<h3 align="center">go-lint</h3>
<p align="center">便捷的 Go 代码检查工具，支持工作区和自动化执行</p>

---

# go-lint

便捷的 Go 代码检查工具，支持工作区和自动化 golangci-lint 执行。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🎯 **便捷的代码检查**: 简单的 golangci-lint 执行，支持工作区发现  
⚡ **项目集支持**: 能够处理单个包、多个根路径和工作区操作  
🔄 **错误处理**: 全面的错误解析和调试功能  
🌍 **工作区集成**: 自动发现 Go 模块，支持可配置过滤选项  
📋 **详细报告**: 彩色输出，包含文件位置和问题描述

## 安装

```bash
go install github.com/go-mate/go-lint/cmd/go-lint@latest
```

## 使用方法

### 基础用法

#### 单项目检查
```bash
cd project-path && go-lint
```
分析并报告当前 Go 项目的代码检查问题。

#### 项目集检查
```bash
cd workspace-path && go-lint
```
自动发现并分析工作区中的 Go 项目集。

### 高级用法

#### 调试模式
```bash
# 启用调试模式获得详细输出
cd project-path && go-lint run --debug

# 显式调试控制
cd project-path && go-lint run --debug=1
cd project-path && go-lint run --debug=0
```

#### 工作区配置
工具使用可配置选项自动发现 Go 模块：
- 包含当前包: ✅ 默认启用
- 包含子模块: ✅ 默认启用  
- 排除非 Go 项目: ✅ 智能过滤
- 调试输出: 🔧 可配置

### 执行模式

#### 1. 单包模式
对当前包执行检查：
```bash
go-lint
```

#### 2. 多根模式  
对多个根路径执行检查：
```bash
go-lint run
```

#### 3. 工作区模式
对整个工作区执行检查：
```bash
go-lint run --debug=1
```

## 技术架构

### 核心组件

1. **统一接口包** (`golint/`): 提供全面的 linting 执行功能
   - 单路径执行 (`Run`)
   - 多根路径处理 (`RootsRun`, `WorksRun`)
   - 带自定义配置的批处理操作 (`BatchRun`)
   - 聚合结果，保持顺序
   - 成功/失败分析

2. **核心执行引擎** (`golangcilint/`): 直接 golangci-lint 集成
   - JSON 结果解析
   - 复杂错误处理
   - 警告处理
   - 详细问题报告

### 错误处理

工具处理多种 golangci-lint 失败场景：

- **完全失败**: 命令执行错误，提供详细诊断
- **部分结果**: 混合成功/警告场景，支持 JSON 提取
- **干净成功**: 标准成功执行，进行结果处理

### 输出格式

- **彩色控制台输出**: 用户友好的视觉反馈
- **详细问题报告**: 文件位置和描述
- **JSON 调试输出**: 用于集成的原始结果
- **进度跟踪**: 项目集执行进度

## 配置

### 工作区选项
```go
config := workspath.NewOptions().
    WithIncludeCurrentPackage(true).
    WithIncludeSubModules(true).
    WithExcludeNoGo(true).
    WithDebugMode(debugMode)
```

### 执行超时
默认超时：每个模块 5 分钟
```go
result := golint.RootsRun(modulePaths, time.Minute*5)
```

## 集成示例

### 示例 1: 底层 golangcilint.Run 使用

此示例演示底层 golangcilint.Run 函数的直接使用，带有完整的结果验证。

```go
package main

import (
	"time"

	"github.com/go-mate/go-lint/golangcilint"
	"github.com/yyle88/must"
	"github.com/yyle88/must/mustslice"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

func main() {
	projectPath := runpath.PARENT.Path()
	zaplog.SUG.Debugln(projectPath)

	// 使用调试配置执行 golangci-lint 并验证结果
	result := golangcilint.Run(osexec.NewCommandConfig().WithDebug(), projectPath, time.Minute*5)
	must.Done(result.Cause)
	mustslice.Have(result.Output)
	mustslice.None(result.Warnings)
	mustslice.None(result.Result.Issues)
}
```

⬆️ **源码:** [源码](internal/demos/demo1x/main.go)

### 示例 2: BatchRun 带自定义配置

此示例展示如何使用 BatchRun 在多个项目上执行检查，带有自定义执行配置。

```go
package main

import (
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

func main() {
	projectPath := runpath.PARENT.Path()
	zaplog.SUG.Debugln(projectPath)

	// 使用调试配置执行批量 linting 并显示结果
	projectPaths := []string{projectPath}
	result := golint.BatchRun(osexec.NewCommandConfig().WithDebug(), projectPaths, time.Minute*5)
	result.DebugIssues()
}
```

⬆️ **源码:** [源码](internal/demos/demo2x/main.go)

### 示例 3: RootsRun 配合 Cobra 框架

此示例演示与 cobra 命令框架的集成，使用 RootsRun 创建 CLI 工具。

```go
package main

import (
	"time"

	"github.com/go-mate/go-lint/golint"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

func main() {
	projectPath := runpath.PARENT.Path()
	zaplog.SUG.Debugln(projectPath)

	projectPaths := []string{projectPath}

	// 创建根命令并提供使用提示
	var rootCmd = &cobra.Command{
		Use:   "lint",
		Short: "lint",
		Long:  "lint",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("Use 'lint run' to execute linting")
		},
	}

	// 添加执行 RootsRun 的 run 子命令
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "golangci-lint run",
		Long:  "golangci-lint run",
		Run: func(cmd *cobra.Command, args []string) {
			result := golint.RootsRun(projectPaths, time.Minute*5)
			result.DebugIssues()
		},
	})

	must.Done(rootCmd.Execute())
}
```

⬆️ **源码:** [源码](internal/demos/demo3x/main.go)

**支持的 golangci-lint 版本:**
```bash
golangci-lint version
```

输出:
```text
golangci-lint has version 2.0.2 built with go1.25.3 from 2b224c2 on 2025-03-25T20:33:26Z
```

## 使用场景

### 单项目开发
- 快速检查当前项目的代码质量
- 集成到开发工作流中
- 提供详细的问题定位

### 项目集工作区
- 批量检查相关项目集
- 工作区级别的代码质量管控
- 统一的检查标准和报告

### 持续集成
- CI/CD 流水线集成
- 自动化代码质量检查
- 可配置的检查策略

### 开发团队协作
- 团队代码规范统一
- 自动化检查流程
- 详细的问题反馈机制

## 命令示例

### 基础检查
```bash
# 检查当前项目
go-lint

# 检查工作区所有项目
go-lint run

# 带调试信息的检查
go-lint run --debug=1
```

### 工作流集成
```bash
# 开发阶段：详细调试输出
go-lint run --debug=1

# CI/CD 阶段：简洁输出
go-lint run --debug=0

# 快速检查：默认模式
go-lint
```

## 最佳实践

### 开发环境配置
1. 安装 go-lint 到全局环境
2. 配置编辑器集成
3. 设置预提交钩子

### 团队使用规范
1. 统一 golangci-lint 版本
2. 制定检查频率标准
3. 建立问题处理流程

### 性能优化建议
1. 合理配置超时时间
2. 使用工作区模式批量处理
3. 根据项目规模调整并发度

---

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-mate/go-lint.svg?variant=adaptive)](https://starchart.cc/go-mate/go-lint)