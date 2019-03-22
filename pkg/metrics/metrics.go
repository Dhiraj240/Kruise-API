package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace = "deploy_wizard"

	counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "request_count",
		Help:      "request count.",
	}, []string{"app", "name", "method", "state"})

	histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "duration_seconds",
		Help:      "Time taken to execute endpoint.",
	}, []string{"app", "name", "method", "status"})

	// AppsRenderedCount tracks the number of rendered apps
	AppsRenderedCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "apps_rendered_count",
		Help:      "apps rendered count.",
	}, []string{"app"})
)
