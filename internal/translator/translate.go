package translator

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure"
)

type columnCoordinates struct {
	index int
	name  string
}

var (
	url                    = "https://google-translate1.p.rapidapi.com/language/translate/v2"
	firstColumnToTranslate = columnCoordinates{16, "Q"}
	concurrentRequestLimit = 10
)

type Translator struct {
	client infrastructure.TranslationClient
}

func NewTranslator(client infrastructure.TranslationClient) Translator {
	return Translator{client}
}

func (t Translator) Run(
	ctx context.Context, sourceFilename, sheetName, targetFilename string) error {
	source, err := excelize.OpenFile(sourceFilename)
	if err != nil {
		return err
	}
	defer func() {
		if err := source.Close(); err != nil {
			log.Printf("can not close the file %s: %v", sourceFilename, err)
		}
	}()
	err = t.translateExcelSheet(ctx, source, sheetName)
	if err != nil {
		return err
	}
	if err := source.SaveAs(targetFilename); err != nil {
		return fmt.Errorf("can not save a file '%v': %w", targetFilename, err)
	}
	return nil
}

func (t Translator) translateExcelSheet(
	ctx context.Context, file *excelize.File, sheetName string) error {
	rows, err := file.Rows(sheetName)
	if err != nil {
		return fmt.Errorf(
			"can not get rows of a sheet '%v': %w", sheetName, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("can not close a sheet '%v'", sheetName)
		}
	}()
	rows.Next()
	firstRow, err := rows.Columns()
	if err != nil {
		return fmt.Errorf(
			"can not read a first row of a sheet '%v': %w", sheetName, err)
	}
	column := columnCoordinates{0, "A"}
	if err := t.translateRow(ctx, 1, column, firstRow, sheetName, file); err != nil {
		return fmt.Errorf("can't translate a first row: %v", err)
	}

	limit := make(chan struct{}, concurrentRequestLimit)
	var waitGroup sync.WaitGroup
	for i := 2; rows.Next(); i++ {
		row, err := rows.Columns()
		if err != nil {
			log.Printf(
				"can not read a row %v of a sheet '%v': %s",
				i,
				sheetName,
				err,
			)
			continue
		}
		if len(row) < firstColumnToTranslate.index+1 {
			log.Printf("a final or unexpected row %v: %v", i, row)
			break
		}

		limit <- struct{}{}
		waitGroup.Add(1)
		go func(rowNumber int) {
			defer waitGroup.Done()
			if err := t.translateRow(
				ctx,
				rowNumber,
				firstColumnToTranslate,
				row,
				sheetName,
				file,
			); err != nil {
				log.Printf("can't translate a row %v: %v", row, err)
			}
			<-limit
		}(i)
	}
	waitGroup.Wait()
	return nil
}

func (t Translator) translateRow(
	ctx context.Context,
	rowNumber int,
	startColumn columnCoordinates,
	rowValues []string,
	sheetName string,
	file *excelize.File,
) error {
	if len(rowValues) < startColumn.index {
		return fmt.Errorf(
			"wrong column index %v, expect less than %v",
			startColumn.index,
			len(rowValues),
		)
	}
	translatedValues := []interface{}{}
	for _, cellValue := range rowValues[startColumn.index:] {
		if num, err := strconv.ParseFloat(cellValue, 32); err == nil {
			translatedValues = append(translatedValues, num)
			continue
		}
		if val, err := strconv.ParseBool(cellValue); err == nil {
			translatedValues = append(translatedValues, val)
			continue
		}
		if cellValue == "" {
			translatedValues = append(translatedValues, cellValue)
			continue
		}
		translation, err := t.getTranslation(ctx, cellValue)
		if err != nil {
			log.Printf("TRANSLATION ERROR: %v", err)
			translatedValues = append(translatedValues, "TRANSLATION ERROR")
			continue
		}
		log.Printf("receive a translation %v", translation)
		translatedValues = append(translatedValues, translation)
	}
	log.Printf("update a row %v: %q\n", rowNumber, translatedValues)
	firstTranslatedCell := fmt.Sprintf("%v%v", startColumn.name, rowNumber)
	if err := file.SetSheetRow(
		sheetName,
		firstTranslatedCell,
		&translatedValues,
	); err != nil {
		log.Printf(
			"can not set a row '%v' for a sheet %v: %s",
			rowNumber,
			sheetName,
			err,
		)
		return fmt.Errorf(
			"can not set a row '%v' for a sheet %v: %w",
			rowNumber,
			sheetName,
			err,
		)
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
