package handlers

import (
	"fmt"
	"reflect"
	"strings"
)

var tagPerLanguage = map[string]string{
	"fi": "nameFi",
	"en": "nameEn",
	"ru": "nameRu",
}
var noDataPerLanguages = map[string]string{
	"fi": "no data",
	"en": "no data",
	"ru": "нет данных",
}

func SerializeIntoMessage(object any, outputLanguage string) (string, error) {
	objectValue := reflect.ValueOf(object)
	if objectValue.Kind() != reflect.Struct {
		return "", fmt.Errorf("not a structure: %v", object)
	}

	var result []string
	t := reflect.TypeOf(object)
	for _, field := range reflect.VisibleFields(t) {
		valueLanguage, ok := field.Tag.Lookup("valueLanguage")
		if !ok {
			return "", fmt.Errorf("no language tag in a field '%s'", field.Name)
		}
		if valueLanguage != outputLanguage && valueLanguage != "all" {
			continue
		}
		featureName, ok := field.Tag.Lookup(tagPerLanguage[outputLanguage])
		if !ok {
			return "", fmt.Errorf("no name tag in a field '%s'", field.Name)
		}
		fieldValue := objectValue.FieldByIndex(field.Index)
		var featureValue string
		switch fieldValue.Kind() {
		case reflect.String:
			featureValue = fieldValue.String()
		case reflect.Int:
			featureValue = fmt.Sprint(fieldValue.Int())
		case reflect.Pointer:
			if fieldValue.IsNil() {
				featureValue = noDataPerLanguages[outputLanguage]
			} else {
				pointerValue := fieldValue.Elem()
				switch pointerValue.Kind() {
				case reflect.String:
					featureValue = pointerValue.String()
				case reflect.Int:
					featureValue = fmt.Sprint(fieldValue.Int())
				}
			}
		default:
			return "", fmt.Errorf("unexpected type of the field %s", field.Name)
		}
		cleanName := strings.ReplaceAll(featureName, "_", " ")
		result = append(result, fmt.Sprintf("%s: %s", cleanName, featureValue))
	}

	return strings.Join(result, "\n"), nil
}