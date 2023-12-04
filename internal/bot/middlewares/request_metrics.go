package middlewares

import (
	"strconv"
	"time"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
)

type BotWithMetrics struct {
	clientName string
	*tgbotapi.BotAPI
	m *metrics.Metrics
}

func NewBotWithMetrics(bot *tgbotapi.BotAPI, m *metrics.Metrics) *BotWithMetrics {
	return &BotWithMetrics{"Telegram", bot, m}
}

func (b *BotWithMetrics) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	start := time.Now()
	msg, err := b.BotAPI.Send(c)
	var isError bool
	if err != nil {
		isError = true
	}
	b.m.RequestDuration.With(
		prometheus.Labels{
			"client": b.clientName,
			"method": "Send",
			"is_error": strconv.FormatBool(isError),
		},
	).Observe(time.Since(start).Seconds())
	return msg, err
}

func (b *BotWithMetrics) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	start := time.Now()
	response, err := b.BotAPI.Request(c)
	var isError bool
	if err != nil {
		isError = true
	}
	b.m.RequestDuration.With(
		prometheus.Labels{
			"client": b.clientName,
			"method": "Request",
			"is_error": strconv.FormatBool(isError),
		},
	).Observe(time.Since(start).Seconds())
	return response, err
}
