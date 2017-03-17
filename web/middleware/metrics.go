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

	httpRequestsDuration = prometheus.NewHistogramVec(
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
	prometheus.MustRegister(httpRequestsDuration)
}


func UpdateMetrics(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {

		// httpsnoop provides convenient access to some exit-side data such as the status code, and a duration
		snoopMetrics := httpsnoop.CaptureMetrics(h, w, req)
		
		normPath := normalizePathForMetric(req.URL.Path)

		statusString := strconv.Itoa(snoopMetrics.Code)

		standardLabels := prometheus.Labels{
				"httpHost": req.Host,
				"httpMethod": req.Method,
				"normalizedPath": normPath,
				"httpStatus": statusString}

		httpRequests.With(standardLabels).Inc()
		httpRequestsDuration.With(standardLabels).Observe( snoopMetrics.Duration.Seconds() )
	}

	return http.HandlerFunc(fn)
}

func normalizePathForMetric(path string) string {
	const placeholder = "/{id}"
	return pathIdPattern.ReplaceAllString(path, placeholder)
}