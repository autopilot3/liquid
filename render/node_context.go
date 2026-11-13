package render

import (
	"github.com/autopilot3/liquid/expressions"
)

// nodeContext provides the evaluation context for rendering the AST.
//
// This type has a clumsy name so that render.Context, in the public API, can
// have a clean name that doesn't stutter.
type nodeContext struct {
	bindings          map[string]interface{}
	config            Config
	findVariablesOnly bool
}

// newNodeContext creates a new evaluation context.
func newNodeContext(scope map[string]interface{}, c Config) nodeContext {
	return nodeContext{
		bindings: copyBindings(scope),
		config:   c,
	}
}

// copyBindings copies only the top level. Rendering only writes top-level keys
// (assign, capture, loop and assigned vars), and a deep copy costs the size of
// every nested value (e.g. a large json binding) on each render.
func copyBindings(scope map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(scope))
	for k, v := range scope {
		out[k] = v
	}
	return out
}

func newFindVariablesNodeContext(c Config) nodeContext {
	return nodeContext{
		bindings:          make(map[string]interface{}),
		config:            c,
		findVariablesOnly: true,
	}
}

// Evaluate evaluates an expression within the template context.
func (c nodeContext) Evaluate(expr expressions.Expression) (out interface{}, err error) {
	if c.findVariablesOnly {
		return expr.Evaluate(expressions.NewVariablesContext(c.bindings, c.config.Config.Config))
	}
	return expr.Evaluate(expressions.NewContext(c.bindings, c.config.Config.Config))
}
