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
<p align="center">Convenient Go linting package with workspace support and automated execution</p>

---

# go-lint

Convenient Go linting package with workspace support and automated golangci-lint execution.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🎯 **Convenient Linting**: Simple golangci-lint execution with workspace detection  
⚡ **Project Set Support**: Capable of handling single packages, multiple root paths, and workspace operations  
🔄 **Error Handling**: Comprehensive error parsing and debugging capabilities  
🌍 **Workspace Integration**: Auto discovers Go modules with configurable filtering options  
📋 **Detailed Reporting**: Colorized output with file positions and issue descriptions

## Installation

```bash
go install github.com/go-mate/go-lint/cmd/go-lint@latest
```

## Usage

### Basic Usage

#### Single Project Linting
```bash
cd project-path && go-lint
```
Analyze and report lint issues in the current Go project.

#### Project Set
```bash
cd workspace-path && go-lint
```
Auto detect and analyze Go project set in the workspace.

### Advanced Usage

#### Debug Mode
```bash
# Enable debug mode with detailed output
cd project-path && go-lint run --debug

# Explicit debug settings
cd project-path && go-lint run --debug=1
cd project-path && go-lint run --debug=0
```

#### Workspace Configuration
This package auto discovers Go modules using configurable options:
- Include current package: ✅ Enabled by default
- Include submodules: ✅ Enabled by default  
- Exclude non-Go projects: ✅ Smart filtering
- Debug output: 🔧 Configurable

### Execution Modes

#### 1. Single Package Mode
Execute linting on the current package:
```bash
go-lint
```

#### 2. Multi-Root Mode  
Execute linting on multiple root paths:
```bash
go-lint run
```

#### 3. Workspace Mode
Execute linting across entire workspace:
```bash
go-lint run --debug=1
```

## Technical Architecture

### Core Components

1. **Unified Interface Package** (`golint/`): Provides comprehensive linting execution
   - Single path execution (`Run`)
   - Multi-root path processing (`RootsRun`, `WorksRun`)
   - Batch operations with custom configuration (`BatchRun`)
   - Aggregated results with sequence preservation
   - Success/failure analysis

2. **Core Execution Engine** (`golangcilint/`): Direct golangci-lint integration
   - JSON result parsing
   - Complex error handling
   - Warning processing
   - Detailed issue reporting

### Error Handling

This package handles multiple golangci-lint failure scenarios:

- **Complete Failures**: Command execution errors with detailed diagnostics
- **Partial Results**: Mixed success/warning scenarios with JSON extraction
- **Clean Success**: Standard success with result processing

### Output Formats

- **Colorized Console Output**: Intuitive visual feedback
- **Detailed Issue Reports**: File positions and descriptions
- **JSON Debug Output**: Raw results for integration
- **Progress Tracking**: Project set execution progress

## Configuration

### Workspace Options
```go
config := workspath.NewOptions().
    WithIncludeCurrentPackage(true).
    WithIncludeSubModules(true).
    WithExcludeNoGo(true).
    WithDebugMode(debugMode)
```

### Execution Timeout
Default timeout: 5 minutes each module
```go
result := golint.RootsRun(modulePaths, time.Minute*5)
```

## Integration Examples

### Example 1: Base-Level golangcilint.Run Usage

This example demonstrates direct usage of the base-level golangcilint.Run function with comprehensive result validation.

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

	// Execute golangci-lint with debug config and validate results
	result := golangcilint.Run(osexec.NewCommandConfig().WithDebug(), projectPath, time.Minute*5)
	must.Done(result.Cause)
	mustslice.Have(result.Output)
	mustslice.None(result.Warnings)
	mustslice.None(result.Result.Issues)
}
```

⬆️ **Source:** [Source](internal/demos/demo1x/main.go)

### Example 2: BatchRun with Custom Configuration

This example shows how to use BatchRun to execute linting on multiple projects with custom execution configuration.

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

	// Execute batch linting with debug config and display results
	projectPaths := []string{projectPath}
	result := golint.BatchRun(osexec.NewCommandConfig().WithDebug(), projectPaths, time.Minute*5)
	result.DebugIssues()
}
```

⬆️ **Source:** [Source](internal/demos/demo2x/main.go)

### Example 3: RootsRun with Cobra Framework

This example demonstrates integration with the cobra command framework to create a CLI application using RootsRun.

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

	// Create root command with usage hint
	var rootCmd = &cobra.Command{
		Use:   "lint",
		Short: "lint",
		Long:  "lint",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("Use 'lint run' to execute linting")
		},
	}

	// Add run subcommand that executes RootsRun
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

⬆️ **Source:** [Source](internal/demos/demo3x/main.go)

**Supported golangci-lint version:**
```bash
golangci-lint version
```

Output:
```text
golangci-lint has version 2.0.2 built with go1.25.3 from 2b224c2 on 2025-03-25T20:33:26Z
```

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/go-lint.svg?variant=adaptive)](https://starchart.cc/go-mate/go-lint)
