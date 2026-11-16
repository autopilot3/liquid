package render

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyBindings(t *testing.T) {
	scope := map[string]any{"key": "value"}

	cp := copyBindings(scope)
	cp["key"] = "mutated"
	cp["added"] = true

	require.Equal(t, map[string]any{"key": "value"}, scope)
}
