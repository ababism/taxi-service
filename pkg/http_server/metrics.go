package http_server

import (
	m "gitlab.com/ArtemFed/mts-final-taxi/pkg/metrics"
)

const metricNameSpace = "http_server"

func newMetrics() *metrics {
	req := m.GetOrRegisterCounterVec(m.CounterOpts{
		Namespace:   metricNameSpace,
		Name:        "taxi_http_server_requests_count",
		Description: "Counter of requests received by HTTP server",
	}, []string{"method", "http_code", "route"})

	return &metrics{
		req: req,
	}
}

type metrics struct {
	req *m.CounterVec
}

func (m *metrics) observe(method string, httpCode string, route string) {
	m.req.WithLabelValues(method, httpCode, route).Inc()
}
