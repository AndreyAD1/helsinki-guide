package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	ChatUpdates prometheus.Counter
	UnexpectedUpdates *prometheus.CounterVec
}

func NewMetrics(registerer prometheus.Registerer) *Metrics {
	metrics := Metrics{
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "helsinki_guide",
			Name:      "chat_updates",
			Help:      "number of received chat updates",
		}),
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "helsinki_guide",
			Name:      "unexpected_updates",
			Help:      "number of unexpected chat updates",
		}, []string{"error"}),
	}
	registerer.MustRegister(metrics.ChatUpdates, metrics.UnexpectedUpdates)
	return &metrics
}
