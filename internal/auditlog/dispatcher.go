package auditlog

import (
	"context"
	"encoding/json"
	"fmt"
	"governo/internal/rules"
	"slices"
	"sync"
)

type Dispatcher struct {
	rules []rules.Rule
	mu    sync.RWMutex
}

func NewDispatcher(rules []rules.Rule) *Dispatcher {
	return &Dispatcher{rules: rules}
}

func (d *Dispatcher) UpdateRules(rules []rules.Rule) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.rules = rules
}

func isCloudAssetInventoryEvent(attrs map[string]any) bool {
	_, ok := attrs["asset"].(map[string]any)
	return ok
}

func evaluateAndExecute(ctx context.Context, r rules.Rule, resource string, attrs map[string]any) {
	ok, err := rules.Evaluate(r.Condition, attrs)
	if err != nil {
		fmt.Printf("Error evaluating rule '%s': %v\n", r.Name, err)
		return
	}
	if ok {
		_ = rules.ExecuteActions(ctx, r.Actions, resource, attrs)
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, data []byte) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var entry map[string]any
	if err := json.Unmarshal(data, &entry); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	attrs := entry

	// metrics collectors
	rules.RunMetrics(ctx, attrs)

	// evaluate rules
	for _, r := range d.rules {
		if isCloudAssetInventoryEvent(attrs) {
			// This is an asset event, so only process rules with AssetType.
			if len(r.AssetType) > 0 {
				asset, _ := entry["asset"].(map[string]any)
				assetType, _ := asset["assetType"].(string)
				resource, _ := asset["name"].(string)

				if contains(r.AssetType, assetType) {
					evaluateAndExecute(ctx, r, resource, attrs)
				}
			}
		} else {
			// This is an audit log event, so only process rules with Service or Methods.
			if r.Service != "" || len(r.Methods) > 0 {
				payload, ok := entry["protoPayload"].(map[string]any)
				if !ok {
					continue
				}
				svc, _ := payload["serviceName"].(string)
				method, _ := payload["methodName"].(string)
				resource, _ := payload["resourceName"].(string)

				if (r.Service == "" || r.Service == svc) && (len(r.Methods) == 0 || contains(r.Methods, method)) {
					evaluateAndExecute(ctx, r, resource, attrs)
				}
			}
		}
	}
	return nil
}

func contains(slice []string, s string) bool {
	return slices.Contains(slice, s)
}
