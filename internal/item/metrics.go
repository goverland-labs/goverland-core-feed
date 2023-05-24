package item

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/goverland-labs/feed/internal/metrics"
)

var metricHandleHistogram = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "feed_item",
		Name:      "handle_duration_seconds",
		Help:      "Handle feed item event duration seconds",
		Buckets:   []float64{.001, .005, .01, .025, .05, .1, .5, 1, 2.5, 5, 10},
	}, []string{"type", "error"},
)