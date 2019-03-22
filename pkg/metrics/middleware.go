package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metricResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newMetricResponseWriter(w http.ResponseWriter) *metricResponseWriter {
	return &metricResponseWriter{w, http.StatusOK}
}

func (lrw *metricResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// SetupHandler handler metrics
func SetupHandler(handler http.Handler, app string) http.Handler {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newMetricResponseWriter(w)
		handler.ServeHTTP(lrw, r)
		statusCode := lrw.statusCode
		duration := time.Since(start)
		histogram.WithLabelValues(app, r.URL.String(), r.Method, fmt.Sprintf("%d", statusCode)).Observe(duration.Seconds())
		counter.WithLabelValues(app, r.URL.String(), r.Method, fmt.Sprintf("%d", statusCode)).Inc()
	})

	prometheus.Register(histogram)
	prometheus.Register(counter)
	prometheus.Register(AppsRenderedCount)
	return h
}
