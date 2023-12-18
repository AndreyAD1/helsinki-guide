package clients

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetResponseWithRetry(client *http.Client, request *http.Request) ([]byte, error) {
	for {
		response, err := client.Do(request)
		if err != nil {
			err = fmt.Errorf("can not send a request: %v: %w", request, err)
			return nil, err
		}
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			err = fmt.Errorf(
				"can not read a response body for a request: %v: %w",
				request,
				err,
			)
			return nil, err
		}
		if response.StatusCode >= 500 {
			log.Printf("receive 500 for a request %v: %s", request, responseBody)
			time.Sleep(time.Second * 2)
			continue
		}

		if response.StatusCode != http.StatusOK {
			err = fmt.Errorf(
				"receive an error status code '%v': %v: %v",
				response.StatusCode,
				request,
				string(responseBody),
			)
			return nil, err
		}
		return responseBody, nil
	}
}
