package liquid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderDoesNotChangeBindings(t *testing.T) {
	nested := map[string]any{"name": "a"}
	items := []any{"x", "y"}
	bindings := map[string]any{"nested": nested, "items": items}

	tpl, err := NewEngine().ParseString(`{% assign nested = "b" %}{% for i in items %}{{ i }}{% endfor %}{% capture c %}z{% endcapture %}{{ nested }}{{ c }}`)
	require.NoError(t, err)
	out, err := tpl.RenderString(bindings)
	require.NoError(t, err)

	require.Equal(t, "xybz", out)
	require.Equal(t, map[string]any{"nested": nested, "items": items}, bindings)
	require.Equal(t, map[string]any{"name": "a"}, nested)
	require.Equal(t, []any{"x", "y"}, items)
}

func BenchmarkRenderLargeBindings(b *testing.B) {
	products := make([]any, 5000)
	for i := range products {
		products[i] = map[string]any{"id": i, "title": "title", "price": "29.90", "tags": []any{"a", "b"}}
	}
	bindings := map[string]any{"name": "a", "json": map[string]any{"products": products}}
	tpl, err := NewEngine().ParseString(`Hi {{ name }}`)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := tpl.RenderString(bindings); err != nil {
			b.Fatal(err)
		}
	}
}
