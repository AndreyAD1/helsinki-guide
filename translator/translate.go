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

func (t Translator) Run(ctx context.Context) error {
	filename := "input_dataset.xlsx"
	source, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := source.Close(); err != nil {
			log.Printf("can not close the file %s: %v", filename, err)
		}
	}()
	translated, err := t.getTranslatedFile(ctx, source)
	if err != nil {
		return err
	}
	if err := translated.SaveAs("translated.xlsx"); err != nil {
		return err
	}
	return nil
}

func (t Translator) getTranslatedFile(ctx context.Context, file *excelize.File) (*excelize.File, error) {
	for _, cell := range []string{"U2"} {
		finnishText, err := file.GetCellValue("Lauttasaari", cell)
		if err != nil {
			log.Printf("can not read the cell T2: %v", err)
			continue
		}
		fmt.Println(finnishText)
		newCtx, cancel := context.WithTimeout(ctx, time.Second * 10)
		defer cancel()
		englishText, err := t.client.GetTranslation(newCtx, "fi", "en", finnishText)
		if err != nil {
			log.Printf("can not translate %v: %v", finnishText, err)
			continue
		}
		file.SetCellValue("Lauttasaari", cell, englishText)
		fmt.Println(englishText)
	}
	return file, nil
}
