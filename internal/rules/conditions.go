package rules

import (
	"fmt"

	"github.com/google/cel-go/cel"
)

// Evaluate evaluates a CEL expression against a map of metrics.
func Evaluate(condition string, payload map[string]any) (bool, error) {
	env, err := cel.NewEnv()
	if err != nil {
		return false, fmt.Errorf("failed to create CEL environment: %w", err)
	}

	ast, issues := env.Parse(condition)
	if issues != nil && issues.Err() != nil {
		return false, fmt.Errorf("failed to parse CEL expression: %w", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return false, fmt.Errorf("failed to create CEL program: %w", err)
	}

	out, _, err := prg.Eval(payload)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate CEL expression: %w", err)
	}

	result, ok := out.Value().(bool)
	if !ok {
		return false, fmt.Errorf("CEL expression did not return a boolean value")
	}

	return result, nil
}
