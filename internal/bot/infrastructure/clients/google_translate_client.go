package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GoogleTranslateClient struct {
	endpoint string
	apiKey   string
}

var translatorURL = "https://google-translate1.p.rapidapi.com/language/translate/v2"

type translatedText struct {
	TranslatedText string `json:"translatedText"`
}
type TranslationResponseBody struct {
	Data struct {
		Translations []struct {
			translatedText
		} `json:"translations"`
	} `json:"data"`
}

func NewGoogleClient(apiKey string) GoogleTranslateClient {
	return GoogleTranslateClient{translatorURL, apiKey}
}

func (c GoogleTranslateClient) GetTranslation(
	ctx context.Context,
	source,
	target,
	text string,
) (string, error) {
	body := fmt.Sprintf("source=%v&target=%v&q=%v", source, target, text)
	bodyReader := strings.NewReader(body)
	errContext := fmt.Sprintf(
		"a request to '%v' with a body %s",
		c.endpoint,
		body,
	)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint,
		bodyReader,
	)
	if err != nil {
		err = fmt.Errorf("can not create a request: %v: %w", errContext, err)
		return "", err
	}
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept-Encoding", "application/gzip")
	request.Header.Add("X-RapidAPI-Key", c.apiKey)
	request.Header.Add("X-RapidAPI-Host", "google-translate1.p.rapidapi.com")

	responseBody, err := GetResponseWithRetry(http.DefaultClient, request)
	if err != nil {
		return "", err
	}
	var parsedResponse TranslationResponseBody
	if err = json.Unmarshal(responseBody, &parsedResponse); err != nil {
		err = fmt.Errorf(
			"receive an unexpected response body: '%v': %w",
			string(responseBody),
			err,
		)
		return "", err
	}
	translatedText := parsedResponse.Data.Translations[0].TranslatedText
	return translatedText, nil
}
