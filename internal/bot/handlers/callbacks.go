package handlers

import (
	c "context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func (h HandlerContainer) next(ctx c.Context, query *tgbotapi.CallbackQuery) error {
	// Telegram asks a bot server to explicitly answer every callback call
	defer func() {
		callbackAnswer := tgbotapi.NewCallback(query.ID, "")
		_, err := h.bot.Request(callbackAnswer)
		if err != nil {
			slog.WarnContext(
				ctx,
				fmt.Sprintf("could not answer to a callback %v", query.ID),
				slog.Any(logger.ErrorKey, err),
			)
		}
	}()

	message := query.Message
	if message == nil {
		errMsg := fmt.Sprintf("a callback has no message %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "a callback has no message"},
		).Inc()
		return nil
	}
	msgID := query.Message.MessageID
	chat := query.Message.Chat
	if chat == nil {
		errMsg := fmt.Sprintf("a callback has no chat %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "a callback has no chat"},
		).Inc()
		return nil
	}
	var button NextButton
	if err := json.Unmarshal([]byte(query.Data), &button); err != nil {
		logMsg := fmt.Sprintf(
			"unexpected callback data %v from a message %v and the chat %v",
			query.Data,
			msgID,
			chat.ID,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "unexpected callback data"},
		).Inc()
		return nil
	}
	// I need to extract an address from a message text
	//  instead of using query data because the Telegram API specifies that
	//  query data should be less than 64 bytes.
	firstRow, _, found := strings.Cut(query.Message.Text, "\n")
	logMsg := fmt.Sprintf(unexpectedTextTmpl, query.Message.Text, msgID, chat.ID)
	if !found {
		slog.ErrorContext(ctx, logMsg)
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "unexpected callback message"},
		).Inc()
		return nil
	}
	_, address, found := strings.Cut(firstRow, ":")
	if !found {
		slog.ErrorContext(ctx, logMsg)
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "unexpected callback message"},
		).Inc()
		return nil
	}
	address = strings.TrimSpace(address)
	if err := h.returnAddresses(ctx, chat.ID, address, button.Limit, button.Offset); err != nil {
		return err
	}
	editedMessage := tgbotapi.NewEditMessageReplyMarkup(
		chat.ID,
		msgID,
		tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
		},
	)
	_, err := h.bot.Send(editedMessage)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not edit a message %v: %v", chat.ID, msgID),
			slog.Any(logger.ErrorKey, err),
		)
	}
	return err
}

func (h HandlerContainer) language(ctx c.Context, query *tgbotapi.CallbackQuery) error {
	// Telegram asks a bot server to explicitly answer every callback call
	defer func() {
		callbackAnswer := tgbotapi.NewCallback(query.ID, "")
		_, err := h.bot.Request(callbackAnswer)
		if err != nil {
			slog.WarnContext(
				ctx,
				fmt.Sprintf("could not answer to a callback %v", query.ID),
				slog.Any(logger.ErrorKey, err),
			)
		}
	}()

	message := query.Message
	if message == nil {
		errMsg := fmt.Sprintf("a callback has no message %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		h.metrics.UnexpectedLanguageCallback.With(
			prometheus.Labels{"error": "a callback has no message"},
		).Inc()
		return nil
	}
	msgID := query.Message.MessageID
	if query.From == nil {
		errMsg := fmt.Sprintf("a callback has no chat %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		h.metrics.UnexpectedLanguageCallback.With(
			prometheus.Labels{"error": "a callback has no user"},
		).Inc()
		return nil
	}

	chat := query.Message.Chat
	if chat == nil {
		errMsg := fmt.Sprintf("a callback has no chat %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		h.metrics.UnexpectedLanguageCallback.With(
			prometheus.Labels{"error": "a callback has no chat"},
		).Inc()
		return nil
	}
	var button LanguageButton
	if err := json.Unmarshal([]byte(query.Data), &button); err != nil {
		logMsg := fmt.Sprintf(
			"unexpected callback data %v from a message %v and the chat %v",
			query.Data,
			msgID,
			chat.ID,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		h.metrics.UnexpectedLanguageCallback.With(
			prometheus.Labels{"error": "unexpected callback data"},
		).Inc()
		return nil
	}
	if err := h.userService.SetLanguage(
		ctx,
		query.From.ID,
		button.Language,
	); err != nil {
		h.metrics.UnexpectedLanguageCallback.With(
			prometheus.Labels{"error": "internal error"},
		).Inc()
		sendErr := h.SendMessage(ctx, chat.ID, "Internal error", "")
		return errors.Join(sendErr, err)
	}
	approve := fmt.Sprintf(
		"I will return the building information in %s.",
		languageCodes[button.Language],
	)
	editedMessage := tgbotapi.NewEditMessageTextAndMarkup(
		chat.ID,
		msgID,
		approve,
		tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
		},
	)
	_, err := h.bot.Send(editedMessage)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not edit a message %v: %v", chat.ID, msgID),
			slog.Any(logger.ErrorKey, err),
		)
	}
	return err
}
