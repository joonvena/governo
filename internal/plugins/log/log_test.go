package log

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLogAction_Handle(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	logAction := &LogAction{}

	// Test case 1: Simple message without templating
	params1 := map[string]any{"message": "Hello, world!"}
	attrs1 := map[string]any{}
	err := logAction.Handle(context.Background(), "test-resource", params1, attrs1)
	if err != nil {
		t.Errorf("Handle failed: %v", err)
	}

	// Test case 2: Message with templating
	params2 := map[string]any{"message": "Project ID: {{ asset.resource.data.projectId }}"}
	attrs2 := map[string]any{
		"asset": map[string]any{
			"resource": map[string]any{
				"data": map[string]any{
					"projectId": "my-templated-project-id",
				},
			},
		},
	}
	err = logAction.Handle(context.Background(), "test-resource", params2, attrs2)
	if err != nil {
		t.Errorf("Handle failed: %v", err)
	}

	// Test case 3: Message with templating and missing field
	params3 := map[string]any{"message": "Missing field: {{ non.existent.field }}"}
	attrs3 := map[string]any{}
	if err = logAction.Handle(context.Background(), "test-resource", params3, attrs3); err != nil {
		t.Errorf("Handle failed: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	out, _ := io.ReadAll(r)
	output := string(out)

	expectedOutput1 := "[LOG] resource=test-resource msg=Hello, world!\n"
	expectedOutput2 := "[LOG] resource=test-resource msg=Project ID: my-templated-project-id\n"
	expectedOutput3 := "[LOG] resource=test-resource msg=Missing field: {{non.existent.field:ERR}}\n"

	if !strings.Contains(output, expectedOutput1) {
		t.Errorf("Expected output not found for test case 1. Got: %s", output)
	}
	if !strings.Contains(output, expectedOutput2) {
		t.Errorf("Expected output not found for test case 2. Got: %s", output)
	}
	if !strings.Contains(output, expectedOutput3) {
		t.Errorf("Expected output not found for test case 3. Got: %s", output)
	}
}
