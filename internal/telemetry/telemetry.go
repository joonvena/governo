package telemetry

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var Meter metric.Meter

func Init(ctx context.Context) func() {
	exporterType := os.Getenv("METRICS_EXPORTER")
	if exporterType == "" {
		exporterType = "disabled" // default
	}

	var provider *sdkmetric.MeterProvider
	var shutdown func()

	switch exporterType {
	case "prometheus":
		// Prometheus exporter
		exp, err := prometheus.New()
		if err != nil {
			log.Fatalf("failed to initialize Prometheus exporter: %v", err)
		}

		provider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(exp))
		Meter = provider.Meter("governo")

		// Expose metrics endpoint
		go func() {
			log.Println("Prometheus metrics available at :9464/metrics")
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(":9464", nil); err != nil {
				log.Fatalf("failed to start Prometheus HTTP server: %v", err)
			}
		}()

		shutdown = func() { _ = provider.Shutdown(ctx) }

	case "console":
		// Console exporter
		exp, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		if err != nil {
			log.Fatalf("failed to initialize stdout exporter: %v", err)
		}

		provider = sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(
				sdkmetric.NewPeriodicReader(exp, sdkmetric.WithInterval(5*time.Second)),
			),
		)
		Meter = provider.Meter("governo")

		shutdown = func() { _ = provider.Shutdown(ctx) }
	case "disabled":
		log.Println("Telemetry disabled")
		shutdown = func() {}
		return shutdown
	default:
		log.Fatalf("unknown METRICS_EXPORTER: %s", exporterType)
	}

	// Set provider globally
	otel.SetMeterProvider(provider)

	return shutdown
}
