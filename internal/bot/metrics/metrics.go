package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	ChatUpdates            prometheus.Counter
	UnexpectedUpdates      *prometheus.CounterVec
	UnexpectedNextCallback *prometheus.CounterVec
	CommandDuration        *prometheus.HistogramVec
	ButtonDuration         *prometheus.HistogramVec
	RequestDuration        *prometheus.HistogramVec
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
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "helsinki_guide",
			Name:      "unexpected_next_callback",
			Help:      "number of unexpected callbacks",
		}, []string{"error"}),
		prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "helsinki_guide",
			Name:      "command_duration",
			Help:      "Duration of the command processing.",
			Buckets:   []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1},
		}, []string{"command_name"}),
		prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "helsinki_guide",
			Name:      "button_duration",
			Help:      "Duration of the button processing.",
			Buckets:   []float64{0.1, 0.3, 0.5, 0.6, 0.7, 0.8, 0.9, 1},
		}, []string{"button_name"}),
		prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "helsinki_guide",
			Name:      "request_duration",
			Help:      "Duration of the sending request.",
			Buckets:   []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1},
		}, []string{"client", "method", "is_error"}),
	}
	registerer.MustRegister(
		metrics.ChatUpdates,
		metrics.UnexpectedUpdates,
		metrics.UnexpectedNextCallback,
		metrics.CommandDuration,
		metrics.ButtonDuration,
		metrics.RequestDuration,
	)
	return &metrics
}
