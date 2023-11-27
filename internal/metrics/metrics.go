package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	TelegramUpdateCounter prometheus.Counter
}

func NewMetrics(registerer prometheus.Registerer) *Metrics {
	metrics := Metrics{
		prometheus.NewCounter(prometheus.CounterOpts{}),
	}
	return &metrics
}