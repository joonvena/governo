package auditlog

import (
	"context"
	"testing"

	"governo/internal/plugins/log"
	"governo/internal/rules"
)

// MockMetric is a mock implementation of the Metric interface
type MockMetric struct {
	CollectFunc func(ctx context.Context, attrs map[string]any) error
}

func (m *MockMetric) Collect(ctx context.Context, attrs map[string]any) error {
	if m.CollectFunc != nil {
		return m.CollectFunc(ctx, attrs)
	}
	return nil
}

func TestDispatcher_Dispatch(t *testing.T) {
	// Unregister all metrics and actions to have a clean test environment
	rules.UnregisterAllMetrics()
	rules.UnregisterAllActions()

	// Register the actual LogAction for testing templating
	rules.RegisterAction(&log.LogAction{})

	var metricCalled bool
	rules.RegisterMetrics(&MockMetric{
		CollectFunc: func(ctx context.Context, attrs map[string]any) error {
			metricCalled = true
			return nil
		},
	})

	r := []rules.Rule{
		{
			Service: "storage.googleapis.com",
			Methods: []string{"storage.buckets.create"},
			Actions: []rules.Action{
				{
					Type: "log",
					Params: map[string]any{
						"message": "Audit Log: Bucket created in europe-west1",
					},
				},
			},
		},
		{
			Service: "compute.googleapis.com",
			Methods: []string{"v1.compute.instances.insert"},
		},
	}

	dispatcher := NewDispatcher(r)

	// Test case 1: Matching rule
	data := []byte(`{
		"protoPayload": {
			"serviceName": "storage.googleapis.com",
			"methodName": "storage.buckets.create"
		}
	}`)

	err := dispatcher.Dispatch(context.Background(), data)
	if err != nil {
		t.Errorf("Dispatch failed with error: %v", err)
	}

	if !metricCalled {
		t.Error("Metric collector was not called")
	}

	// For now, we can't easily assert the output of LogAction here without capturing stdout.
	// The log_test.go will cover the templating logic.

	// Test case 2: No matching rule
	metricCalled = false
	data = []byte(`{
		"protoPayload": {
			"serviceName": "iam.googleapis.com",
			"methodName": "google.iam.admin.v1.CreateRole"
		}
	}`)

	err = dispatcher.Dispatch(context.Background(), data)
	if err != nil {
		t.Errorf("Dispatch failed with error: %v", err)
	}

	if !metricCalled {
		t.Error("Metric collector was not called")
	}

	// Test case 3: Cloud Asset Inventory event with templating
	metricCalled = false
	data = []byte(`{
		"asset": {
			"name": "//cloudresourcemanager.googleapis.com/projects/123456789012",
			"assetType": "cloudresourcemanager.googleapis.com/Project",
			"resource": {
				"data": {
					"projectId": "my-templated-project-id"
				}
			}
		}
	}`)

	r = append(r, rules.Rule{
		AssetType: []string{"cloudresourcemanager.googleapis.com/Project"},
		Actions: []rules.Action{
			{
				Type: "log",
				Params: map[string]any{
					"message": "Cloud Asset Inventory: New project created with ID {{ asset.resource.data.projectId }}",
				},
			},
		},
	})

	dispatcher = NewDispatcher(r)

	err = dispatcher.Dispatch(context.Background(), data)
	if err != nil {
		t.Errorf("Dispatch failed with error: %v", err)
	}

	if !metricCalled {
		t.Error("Metric collector was not called")
	}
}
