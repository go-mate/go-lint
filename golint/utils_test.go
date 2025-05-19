package golint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetDebugMode(t *testing.T) {
	SetDebugMode(true)

	require.True(t, debugModeOpen)
}
