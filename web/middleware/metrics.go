package middleware

import (
	"net/http"
	"strconv"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/felixge/httpsnoop"
)

const (
	namespace = "microcosm"
	subsystem = "web"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_total",
			Help:      "Number of http requests handled",
		},
		[]string{"httpHost", "httpMethod", "normalizedPath", "httpStatus"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_duration_seconds",
			Help:      "Seconds to complete the transaction",
			Buckets:   prometheus.ExponentialBuckets(0.001, 10, 7), // 1ms - 1000s
		},
		[]string{"httpHost", "httpMethod", "normalizedPath", "httpStatus"},
	)

	pathIdPattern = regexp.MustCompile("/[0-9]+")

)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpRequestDuration)
}


func UpdateMetrics(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		hsMetrics := httpsnoop.CaptureMetrics(h, w, req)

		//h.ServeHTTP(w, req)
		
		normPath := normalizePathForMetric(req.URL.Path)

		statusString := strconv.Itoa(hsMetrics.Code)

		httpRequestDuration.With(
			prometheus.Labels{
				"httpHost": req.Host,
				"httpMethod": req.Method,
				"normalizedPath": normPath,
				"httpStatus": statusString},
		).Observe( hsMetrics.Duration.Seconds() )

		httpRequests.With(
			prometheus.Labels{
				"httpHost": req.Host,
				"httpMethod": req.Method,
				"normalizedPath": normPath,
				"httpStatus": statusString},
		).Inc()
	}

	return http.HandlerFunc(fn)

}

func normalizePathForMetric(path string) string {
	const placeholder = "/{id}"
	return pathIdPattern.ReplaceAllString(path, placeholder)
}