package http_server

import (
	m "gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"time"
)

const metricNameSpace = "http_server"

func newMetrics() *metrics {
	reqTime := m.GetOrRegisterSummaryVec(m.SummaryOpts{
		Namespace:   metricNameSpace,
		Name:        "taxi_http_server_response_time",
		Description: "Average time by percentile on HTTP server response",
		Objectives:  m.DefaultObjectives,
	}, []string{"method", "http_code", "route"})

	req := m.GetOrRegisterCounterVec(m.CounterOpts{
		Namespace:   metricNameSpace,
		Name:        "taxi_http_server_requests_count",
		Description: "Counter of requests received by HTTP server",
	}, []string{"method", "http_code", "route"})

	return &metrics{
		reqTime: reqTime,
		req:     req,
	}
}

type metrics struct {
	reqTime *m.SummaryVec
	req     *m.CounterVec
}

func (m *metrics) observe(method string, httpCode string, route string, t time.Time) {
	m.reqTime.WithLabelValues(method, httpCode, route).Observe(time.Since(t).Seconds())
	m.req.WithLabelValues(method, httpCode, route).Inc()
}
