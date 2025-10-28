package golint

// Global debug mode flag for controlling verbose output across all linting operations
// Controls detailed logging, issue reporting, and execution tracing
//
// 全局调试模式标志，控制所有 linting 操作的详细输出
// 控制详细日志记录、问题报告和执行追踪
var debugModeOpen = false

// SetDebugMode enables or disables global debug mode for all linting operations
// When enabled, provides detailed execution logs and comprehensive issue reporting
//
// SetDebugMode 启用或禁用所有 linting 操作的全局调试模式
// 启用时提供详细的执行日志和全面的问题报告
func SetDebugMode(enable bool) {
	debugModeOpen = enable
}
