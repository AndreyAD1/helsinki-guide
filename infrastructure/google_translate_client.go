package infrastructure

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)


type GoogleTranslateClient struct {
	endpoint string
	apiKey string
}

var url = "https://google-translate1.p.rapidapi.com/language/translate/v2"

func NewGoogleClient(apiKey string) GoogleTranslateClient {
	return GoogleTranslateClient{url, apiKey}
}

func (client GoogleTranslateClient) GetTranslation(
	ctx context.Context,
	source,
	target,
	text string,
) (string, error) {
	body := fmt.Sprintf("source=%v&target=%v&q=%v", source, target, text)
	bodyReader := strings.NewReader(body)
	errContext := fmt.Sprintf(
		"a request to '%v' with a body %s", 
		client.endpoint, 
		body,
	)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		client.endpoint, 
		bodyReader,
	)
	if err != nil {
		err = fmt.Errorf("can not create a request: %v: %w", errContext, err)
		return "", err
	}
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept-Encoding", "application/gzip")
	request.Header.Add("X-RapidAPI-Key", client.apiKey)
	request.Header.Add("X-RapidAPI-Host", "google-translate1.p.rapidapi.com")
	
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		err = fmt.Errorf("can not send a request: %v: %w", errContext, err)
		return "", err
	}
	defer response.Body.Close()
	
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf(
			"can not read a response body: %v: %w", 
			errContext, 
			err,
		)
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"receive error status code '%v': %v: %v",
			response.StatusCode,
			errContext,
			string(responseBody),
		)
		return "", err
	}
	return string(responseBody), nil
}