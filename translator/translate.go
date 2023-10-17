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

func (t Translator) Run(ctx context.Context) {
	t.readFile(ctx, "input_dataset.xlsx")
}

func (t Translator) readFile(ctx context.Context, filename string) {
	file, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatalf("can not open a file %s: %v", filename, err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("can not close the file %s: %v", filename, err)
		}
	}()
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
		fmt.Println(englishText)
	}
}
