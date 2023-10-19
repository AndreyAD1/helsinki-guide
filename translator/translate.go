package translator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AndreyAD1/helsinki-guide/infrastructure"
	"github.com/xuri/excelize/v2"
)

var url = "https://google-translate1.p.rapidapi.com/language/translate/v2"

type Translator struct {
	client infrastructure.TranslationClient
}

func NewTranslator(client infrastructure.TranslationClient) Translator {
	return Translator{client}
}

func (t Translator) Run(ctx context.Context, sourceFilename, targetFilename string) error {
	source, err := excelize.OpenFile(sourceFilename)
	if err != nil {
		return err
	}
	defer func() {
		if err := source.Close(); err != nil {
			log.Printf("can not close the file %s: %v", sourceFilename, err)
		}
	}()
	err = t.getTranslatedFile(ctx, source)
	if err != nil {
		return err
	}
	if err := source.SaveAs(targetFilename); err != nil {
		return fmt.Errorf("can not save a file '%v': %w", targetFilename, err)
	}
	return nil
}

func (t Translator) getTranslatedFile(ctx context.Context, file *excelize.File) error {
	for _, sheetName := range file.GetSheetList() {
		rows, err := file.Rows(sheetName)
		if err != nil {
			return fmt.Errorf("can not get rows for a sheet '%v': %w", sheetName, err)
		}
		defer func() {
			if err := rows.Close(); err != nil {
				log.Printf("can not close a sheet '%v'", sheetName)
			}
		}()
		rows.Next()
		firstRow, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("can not read a first row of a sheet '%v': %w", sheetName, err)
		}
		translatedValues := []string{}
		for _, cellValue := range firstRow {
			if cellValue == "" {
				translatedValues = append(translatedValues, "")
				continue
			}
			translation, err := t.getTranslation(ctx, cellValue)
			if err != nil {
				log.Printf("a translation error: %v", err)
				continue
			}
			log.Printf("receive a translations %v", translation)
			translatedValues = append(translatedValues, translation)
		}
		log.Printf("update a first row: %q\n", translatedValues)
		if err = file.SetSheetRow(sheetName, "A1", &translatedValues); err != nil {
			return fmt.Errorf(
				"can not set a new first row '%v' for a sheet %v: %w",
				translatedValues,
				sheetName,
				err,
			)
		}
	}
	return nil
}

func (t Translator) getTranslation(ctx context.Context, text string) (string, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	englishText, err := t.client.GetTranslation(newCtx, "fi", "en", text)
	if err != nil {
		err := fmt.Errorf("can not translate %v: %v", text, err)
		return "", err
	}
	return englishText, nil
}
