package rules

import "context"

type MetricsCollector interface {
	Collect(ctx context.Context, attrs map[string]any) error
}

var metricsCollectors []MetricsCollector

func RegisterMetrics(c MetricsCollector) {
	metricsCollectors = append(metricsCollectors, c)
}

func RunMetrics(ctx context.Context, attrs map[string]any) {
	for _, c := range metricsCollectors {
		_ = c.Collect(ctx, attrs)
	}
}

func UnregisterAllMetrics() {
	metricsCollectors = []MetricsCollector{}
}
