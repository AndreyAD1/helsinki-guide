package handlers

import (
	c "context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Telegram requires bots to explicitly answer every callback call
func (h HandlerContainer) getCallbackAnswerFunc(ctx c.Context, queryID string) func() {
	return func() {
		callbackAnswer := tgbotapi.NewCallback(queryID, "")
		_, err := h.bot.Request(callbackAnswer)
		if err != nil {
			slog.WarnContext(
				ctx,
				fmt.Sprintf("could not answer to a callback %v", queryID),
				slog.Any(logger.ErrorKey, err),
			)
		}
	}
}

func (h HandlerContainer) next(ctx c.Context, query *tgbotapi.CallbackQuery) error {
	defer h.getCallbackAnswerFunc(ctx, query.ID)()

	message := query.Message
	if message == nil {
		err := fmt.Errorf("a callback has no message %v", query.ID)
		slog.WarnContext(ctx, err.Error())
		return errors.Join(err, ErrUnexpectedCallback)
	}
	msgID := query.Message.MessageID
	chat := query.Message.Chat
	if chat == nil {
		err := fmt.Errorf("a callback has no chat %v", query.ID)
		slog.WarnContext(ctx, err.Error())
		return errors.Join(err, ErrUnexpectedCallback)
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
		return errors.Join(err, ErrUnexpectedCallback)
	}
	// I need to extract an address from a message text
	//  instead of using query data because the Telegram API specifies that
	//  query data should be less than 64 bytes.
	firstRow, _, found := strings.Cut(query.Message.Text, "\n")
	logMsg := fmt.Sprintf(unexpectedTextTmpl, query.Message.Text, msgID, chat.ID)
	if !found {
		slog.ErrorContext(ctx, logMsg)
		return fmt.Errorf("%v: %w", logMsg, ErrUnexpectedCallback)
	}
	_, address, found := strings.Cut(firstRow, ":")
	if !found {
		slog.ErrorContext(ctx, logMsg)
		return fmt.Errorf("%v: %w", logMsg, ErrUnexpectedCallback)
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
	defer h.getCallbackAnswerFunc(ctx, query.ID)()
	message := query.Message
	if message == nil {
		err := fmt.Errorf("a callbacbuttonk has no message %v", query.ID)
		slog.WarnContext(ctx, err.Error())
		return errors.Join(err, ErrUnexpectedCallback)
	}
	if query.From == nil {
		err := fmt.Errorf("a callback has no sender %v", query.ID)
		slog.WarnContext(ctx, err.Error())
		return errors.Join(err, ErrUnexpectedCallback)
	}
	chat := query.Message.Chat
	if chat == nil {
		errMsg := fmt.Sprintf("a callback has no chat %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		return fmt.Errorf("%v: %w", errMsg, ErrUnexpectedCallback)
	}
	msgID := query.Message.MessageID
	var button LanguageButton
	if err := json.Unmarshal([]byte(query.Data), &button); err != nil {
		logMsg := fmt.Sprintf(
			"unexpected callback data %v from a message %v and the chat %v",
			query.Data,
			msgID,
			chat.ID,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return fmt.Errorf("%v: %w", logMsg, ErrUnexpectedCallback)
	}
	language, ok := services.GetLanguagePerCode(button.Language)
	if !ok {
		err := fmt.Errorf("unexpected button language '%v': %v", button, msgID)
		slog.ErrorContext(
			ctx,
			err.Error(),
		)
		sendErr := h.SendMessage(ctx, chat.ID, "Internal error", "")
		return errors.Join(sendErr, err)
	}
	if err := h.userService.SetLanguage(
		ctx,
		query.From.ID,
		language,
	); err != nil {
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

func (h HandlerContainer) building(ctx c.Context, query *tgbotapi.CallbackQuery) error {
	defer h.getCallbackAnswerFunc(ctx, query.ID)()
	message := query.Message
	if message == nil {
		err := fmt.Errorf("a callback has no message %v", query.ID)
		slog.WarnContext(ctx, err.Error())
		return errors.Join(err, ErrUnexpectedCallback)
	}
	if query.From == nil {
		err := fmt.Errorf("a callback has no sender %v", query.ID)
		slog.WarnContext(ctx, err.Error())
		return errors.Join(err, ErrUnexpectedCallback)
	}
	chat := query.Message.Chat
	if chat == nil {
		errMsg := fmt.Sprintf("a callback has no chat %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		return fmt.Errorf("%v: %w", errMsg, ErrUnexpectedCallback)
	}
	msgID := query.Message.MessageID
	_, err := h.userService.GetPreferredLanguage(ctx, query.From.ID)
	if err != nil {
		sendErr := h.SendMessage(ctx, chat.ID, "Internal error", "")
		return errors.Join(sendErr, err)
	}
	var button BuildingButton
	if err := json.Unmarshal([]byte(query.Data), &button); err != nil {
		logMsg := fmt.Sprintf(
			"unexpected callback data %v from a message %v and the chat %v",
			query.Data,
			msgID,
			chat.ID,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return fmt.Errorf("%v: %w", logMsg, ErrUnexpectedCallback)
	}
	building, err := h.buildingService.GetBuildingByID(ctx, button.ID)
	userLanguage := services.English
	if user := message.From; user != nil {
		switch user.LanguageCode {
		case "fi":
			userLanguage = services.Finnish
		case "ru":
			userLanguage = services.Russian
		}
		preferredLanguage, err := h.userService.GetPreferredLanguage(
			ctx,
			user.ID,
		)
		if err == nil && preferredLanguage != nil {
			userLanguage = *preferredLanguage
		}
	}
	serializedItem, err := SerializeIntoMessage(*building, userLanguage)
	if err != nil {
		slog.ErrorContext(
			ctx,
			fmt.Sprintf("can not serialize a building '%v'", button.ID),
			slog.Any(logger.ErrorKey, err),
		)
	}
	return h.SendMessage(ctx, message.Chat.ID, serializedItem, tgbotapi.ModeHTML)
}
