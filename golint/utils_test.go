// Package golint: Tests utilities and debug mode settings
// Validates debug mode flag manipulation
//
// golint: 测试实用程序和调试模式设置
// 验证调试模式标志操作
package golint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSetDebugMode validates debug mode flag setting
// Ensures SetDebugMode updates the package-wide debug flag
//
// TestSetDebugMode 验证调试模式标志设置
// 确保 SetDebugMode 更新包级调试标志
func TestSetDebugMode(t *testing.T) {
	SetDebugMode(true)

	require.True(t, debugModeOpen)
}
