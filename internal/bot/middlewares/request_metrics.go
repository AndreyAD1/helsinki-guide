package middlewares

import (
	"strconv"
	"time"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func Duration(
	f func() (interface{}, error), 
	m *metrics.Metrics, 
	clientName, 
	methodName string,
) (interface{}, error) {
	start := time.Now()
	result, err := f()
	var isError bool
	if err != nil {
		isError = true
	}
	m.RequestDuration.With(
		prometheus.Labels{
			"client": clientName,
			"method": methodName,
			"is_error": strconv.FormatBool(isError),
		},
	).Observe(time.Since(start).Seconds())
	return result, err
}
