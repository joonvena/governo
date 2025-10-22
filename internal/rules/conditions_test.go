package rules

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	payload := map[string]any{
		"cost":     150.0,
		"currency": "USD",
		"project":  "gcp-project-123",
		"tags":     []string{"prod", "frontend"},
	}

	testCases := []struct {
		name      string
		condition string
		expected  bool
		expectErr bool
	}{
		{
			name:      "Greater than match",
			condition: "cost > 100.0",
			expected:  true,
		},
		{
			name:      "Greater than mismatch",
			condition: "cost > 200.0",
			expected:  false,
		},
		{
			name:      "Less than match",
			condition: "cost < 200.0",
			expected:  true,
		},
		{
			name:      "Less than mismatch",
			condition: "cost < 100.0",
			expected:  false,
		},
		{
			name:      "Equals match",
			condition: "project == 'gcp-project-123'",
			expected:  true,
		},
		{
			name:      "Equals mismatch",
			condition: "project == 'wrong-project'",
			expected:  false,
		},
		{
			name:      "Logical AND match",
			condition: "cost > 100.0 && currency == 'USD'",
			expected:  true,
		},
		{
			name:      "Logical AND mismatch",
			condition: "cost > 200.0 && currency == 'USD'",
			expected:  false,
		},
		{
			name:      "List membership match",
			condition: "'prod' in tags",
			expected:  true,
		},
		{
			name:      "List membership mismatch",
			condition: "'dev' in tags",
			expected:  false,
		},
		{
			name:      "Invalid expression",
			condition: "cost >> 100.0",
			expectErr: true,
		},
		{
			name:      "Non-boolean result",
			condition: "cost + 10.0",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Evaluate(tc.condition, payload)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}
