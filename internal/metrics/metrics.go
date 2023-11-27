package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	ChatUpdates prometheus.Counter
}

func NewMetrics(registerer prometheus.Registerer) *Metrics {
	metrics := Metrics{
		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "helsinki_guide",
			Name:      "chat_updates",
			Help:      "number of received chat updates",
		}),
	}
	registerer.MustRegister(metrics.ChatUpdates)
	return &metrics
}
