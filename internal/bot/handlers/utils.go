package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	s "github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tagPerLanguage = map[s.Language]string{
	s.Finnish: "nameFi",
	s.English: "nameEn",
	s.Russian: "nameRu",
}
var noDataPerLanguages = map[s.Language]string{
	s.Finnish: "no data",
	s.English: "no data",
	s.Russian: "нет данных",
}
var ErrUnexpectedType error = errors.New("unexpected input type")
var ErrUnexpectedFieldType error = errors.New("unexpected field type")
var ErrNoFieldTag error = errors.New("no expected field tag")
var ErrNoNameTag error = errors.New("no name tag")

func SerializeIntoMessage(object any, outputLanguage s.Language) (string, error) {
	objectValue := reflect.ValueOf(object)
	if objectValue.Kind() != reflect.Struct {
		return "", fmt.Errorf("not a structure: %v: %w", object, ErrUnexpectedType)
	}

	var result []string
	t := reflect.TypeOf(object)
	for _, field := range reflect.VisibleFields(t) {
		valueLanguage, ok := field.Tag.Lookup("valueLanguage")
		if !ok {
			return "", fmt.Errorf(
				"no language tag for the field '%s': %w",
				field.Name,
				ErrNoFieldTag,
			)
		}
		if valueLanguage != string(outputLanguage) && valueLanguage != "all" {
			continue
		}
		featureName, ok := field.Tag.Lookup(tagPerLanguage[outputLanguage])
		if !ok {
			return "", fmt.Errorf(
				"no name tag for the field '%s': %w",
				field.Name,
				ErrNoNameTag,
			)
		}
		fieldValue := objectValue.FieldByIndex(field.Index)
		var featureValue string
		switch fieldValue.Kind() {
		case reflect.String:
			featureValue = fieldValue.String()
		case reflect.Int:
			featureValue = fmt.Sprint(fieldValue.Int())
		case reflect.Slice, reflect.Array:
			items := []string{}
			for i := 0; i < fieldValue.Len(); i++ {
				items = append(items, fieldValue.Index(i).String())
			}
			featureValue = strings.Join(items, ", ")
		case reflect.Pointer:
			if fieldValue.IsNil() {
				featureValue = noDataPerLanguages[outputLanguage]
			} else {
				pointerValue := fieldValue.Elem()
				switch pointerValue.Kind() {
				case reflect.String:
					featureValue = pointerValue.String()
				case reflect.Int:
					featureValue = fmt.Sprint(pointerValue.Int())
				case reflect.Slice, reflect.Array:
					items := []string{}
					for i := 0; i < pointerValue.Len(); i++ {
						items = append(items, pointerValue.Index(i).String())
					}
					featureValue = strings.Join(items, ", ")
				default:
					return "", fmt.Errorf(
						"unexpected type of the field '%s': %w",
						field.Name,
						ErrUnexpectedFieldType,
					)
				}
			}
		default:
			return "", fmt.Errorf(
				"unexpected type of the field '%s': %w",
				field.Name,
				ErrUnexpectedFieldType,
			)
		}
		cleanName := strings.ReplaceAll(featureName, "_", " ")
		result = append(result, fmt.Sprintf("<b>%s:</b> %s", cleanName, featureValue))
	}

	return strings.Join(result, "\n"), nil
}

func getBuildingButtonRows(
	ctx context.Context,
	buildings []services.BuildingPreview,
) ([][]tgbotapi.InlineKeyboardButton, error) {
	keyboardRows := [][]tgbotapi.InlineKeyboardButton{}
	for i, building := range buildings {
		label := fmt.Sprintf(
			lineTemplate,
			i+1,
			building.Address,
			building.Name,
		)
		button := BuildingButton{
			Button{label, BUILDING_BUTTON},
			strconv.FormatInt(building.ID, 10),
		}
		buttonCallbackData, err := json.Marshal(button)
		if err != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("can not create a button %v", button),
				slog.Any(logger.ErrorKey, err),
			)
			return [][]tgbotapi.InlineKeyboardButton{}, err
		}
		buttonData := tgbotapi.NewInlineKeyboardButtonData(
			button.label,
			string(buttonCallbackData),
		)
		buttonRow := tgbotapi.NewInlineKeyboardRow(buttonData)
		keyboardRows = append(keyboardRows, buttonRow)
	}
	return keyboardRows, nil
}
