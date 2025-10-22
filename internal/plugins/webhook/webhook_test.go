package webhook

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookAction_Handle(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tests := []struct {
		name    string
		params  map[string]any
		attrs   map[string]any
		wantErr bool
	}{
		{
			name: "should post to webhook",
			params: map[string]any{
				"url":    server.URL,
				"method": "POST",
				"body":   `{"message": "hello"}`,
			},
			attrs:   map[string]any{},
			wantErr: false,
		},
		{
			name: "should post to webhook with templating",
			params: map[string]any{
				"url":    server.URL,
				"method": "POST",
				"body":   `{"message": "hello {{.foo}}"}`,
			},
			attrs: map[string]any{
				"foo": "bar",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebhookAction{}
			if err := w.Handle(context.Background(), "test-resource", tt.params, tt.attrs); (err != nil) != tt.wantErr {
				t.Errorf("WebhookAction.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
