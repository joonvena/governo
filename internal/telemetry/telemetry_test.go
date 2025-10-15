package telemetry

import (
	"context"
	"os"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestInit(t *testing.T) {
	testCases := []struct {
		name          string
		exporterEnv   string
		expectedMeter bool
	}{
		{
			name:          "Console exporter",
			exporterEnv:   "console",
			expectedMeter: true,
		},
		{
			name:          "Prometheus exporter",
			exporterEnv:   "prometheus",
			expectedMeter: true,
		},
		{
			name:          "Default exporter",
			exporterEnv:   "",
			expectedMeter: false,
		},
		{
			name:          "Disabled exporter",
			exporterEnv:   "disabled",
			expectedMeter: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			if err := os.Setenv("METRICS_EXPORTER", tc.exporterEnv); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.Unsetenv("METRICS_EXPORTER"); err != nil {
					t.Fatal(err)
				}
			}()

			// Reset global meter provider
			otel.SetMeterProvider(nil)
			Meter = nil

			// Call Init
			shutdown := Init(context.Background())
			defer shutdown()

			// Check if meter is initialized
			if (Meter != nil) != tc.expectedMeter {
				t.Errorf("Expected meter to be initialized: %v, but got: %v", tc.expectedMeter, Meter != nil)
			}
		})
	}
}
