package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeepCopyBindings(t *testing.T) {
	t.Run("map is copied not aliased", func(t *testing.T) {
		original := map[string]any{"key": "value"}
		scope := map[string]any{"nested": original}

		cp := deepCopyBindings(scope)

		cp["nested"].(map[string]any)["key"] = "mutated"
		require.Equal(t, "value", original["key"], "original map must not be affected by copy mutation")
	})

	t.Run("slice is copied not aliased", func(t *testing.T) {
		original := []any{"a", "b"}
		scope := map[string]any{"items": original}

		cp := deepCopyBindings(scope)

		cp["items"].([]any)[0] = "mutated"
		assert.Equal(t, "a", original[0], "original slice must not be affected by copy mutation")
	})

	t.Run("scalar value passes through unchanged", func(t *testing.T) {
		scope := map[string]any{"n": 42}
		cp := deepCopyBindings(scope)
		assert.Equal(t, 42, cp["n"])
	})
}
