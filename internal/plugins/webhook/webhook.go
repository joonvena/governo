package webhook

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"governo/internal/plugins/common"
)

type WebhookAction struct{}

func (w *WebhookAction) Type() string { return "webhook" }

func (w *WebhookAction) Handle(ctx context.Context, resource string, params map[string]any, attrs map[string]any) error {
	url, _ := params["url"].(string)
	method, _ := params["method"].(string)
	body, _ := params["body"].(string)

	processedBody := common.ProcessTemplate(body, attrs)

	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(processedBody))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close webhook response body: %v", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
