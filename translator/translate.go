package translator

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func translate(ctx context.Context, text, targetLanguage string) (string, error) {
	url := "https://google-translate1.p.rapidapi.com/language/translate/v2"
	body := strings.NewReader()
	request, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		err = fmt.Errorf(
			"can not create a request to '%v' with a body %s: %w", 
			url, 
			body, 
			err,
		)
		return "", err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		err = fmt.Errorf(
			"can not send a request to '%v' with a body %s: %w", 
			url, 
			body, 
			err,
		)
		return "", err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf(
			"can not read a response body for a request to '%v' with a body %s: %w", 
			url, 
			body, 
			err,
		)
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"Receive error status code '%v' for a request to '%v' with a body %s: %v",
			response.StatusCode,
			url,
			body,
			string(responseBody),
		)
		return "", err
	}
	return string(responseBody), nil
}


func readFile(ctx context.Context, filename string) {
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
	for _, cell := range []string{"C2", "T2", "U2"} {
		finnishText, err := file.GetCellValue("Lauttasaari", cell)
		if err != nil {
			log.Printf("can not read the cell T2: %v", err)
			continue
		}
		fmt.Println(finnishText)
		newCtx, cancel := context.WithTimeout(ctx, time.Second * 10)
		defer cancel()
		englishText, err := translate(newCtx, finnishText, "en")
		if err != nil {
			log.Printf("can not translate %v: %v", finnishText, err)
			continue
		}
		fmt.Println(englishText)
	}
}


func Run(ctx context.Context) {
	readFile(ctx, "input_dataset.xlsx")
}