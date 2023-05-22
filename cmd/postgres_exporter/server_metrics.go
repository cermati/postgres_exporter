package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ServerMetrics holds metrics about a Server. Use NewServerMetrics() to create
// a new instance.
//
// Metrics collected by Collect() will be labeled by datname="<databaseName>"
// and serverConstLabels passed to NewServerMetrics().
type ServerMetrics struct {
	// Last query duration (in seconds) for every metric query (including default
	// and settings metrics)
	MetricQueryDuration *prometheus.GaugeVec

	// Counter of query errors for every metric query (including default and
	// settings metrics)
	MetricQueryErrorTotal *prometheus.CounterVec
}

// NewServerMetrics creates a new ServerMetrics
func NewServerMetrics(databaseName string, serverConstLabels prometheus.Labels) (*ServerMetrics, error) {
	constLabels := make(prometheus.Labels)
	for k, v := range serverConstLabels {
		constLabels[k] = v
	}

	constLabels["datname"] = databaseName

	sm := &ServerMetrics{
		MetricQueryDuration: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Subsystem:   exporter,
				Name:        "metric_query_last_duration_seconds",
				Help:        "Duration of the last metric query",
				ConstLabels: constLabels,
			},
			[]string{"query"},
		),

		MetricQueryErrorTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   exporter,
				Name:        "metric_query_errors_total",
				Help:        "Number of metric query execution errors",
				ConstLabels: constLabels,
			},
			[]string{"query"},
		),
	}

	return sm, nil
}

// Collect collects metrics for the Prometheus client
func (sm *ServerMetrics) Collect(ch chan<- prometheus.Metric) {
	sm.MetricQueryDuration.Collect(ch)
	sm.MetricQueryErrorTotal.Collect(ch)
}

// RecordMetricQueryExecution records the execution duration of a metric query
// as well as whether the query failed or not
func (sm *ServerMetrics) RecordMetricQueryExecution(
	queryName string,
	duration time.Duration,
	queryErr error,
) {
	sm.MetricQueryDuration.WithLabelValues(queryName).Set(duration.Seconds())

	if queryErr != nil {
		sm.MetricQueryErrorTotal.WithLabelValues(queryName).Inc()
	} else {
		// We do this so the metric will be initialized even when all previous query
		// executions didn't fail
		sm.MetricQueryErrorTotal.WithLabelValues(queryName).Add(0)
	}
}
