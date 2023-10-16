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
var url = "https://google-translate1.p.rapidapi.com/language/translate/v2"


func translate(ctx context.Context, apiKey, text, targetLanguage string) (string, error) {
	body := fmt.Sprintf("source=fi&target=en&q=%v", text)
	body_reader := strings.NewReader(body)
	request, err := http.NewRequestWithContext(ctx, "POST", url, body_reader)
	if err != nil {
		err = fmt.Errorf(
			"can not create a request to '%v' with a body %s: %w", 
			url, 
			body, 
			err,
		)
		return "", err
	}
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept-Encoding", "application/gzip")
	request.Header.Add("X-RapidAPI-Key", apiKey)
	request.Header.Add("X-RapidAPI-Host", "google-translate1.p.rapidapi.com")
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
			"receive error status code '%v' for a request to '%v' with a body %s: %v",
			response.StatusCode,
			url,
			body,
			string(responseBody),
		)
		return "", err
	}
	return string(responseBody), nil
}


func readFile(ctx context.Context, filename, apiKey string) {
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
		englishText, err := translate(newCtx, apiKey, finnishText, "en")
		if err != nil {
			log.Printf("can not translate %v: %v", finnishText, err)
			continue
		}
		fmt.Println(englishText)
	}
}


func Run(ctx context.Context, apiKey string) {
	readFile(ctx, "input_dataset.xlsx", apiKey)
}