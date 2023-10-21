package infrastructure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGoogleTranslateClient_GetResponseWithRetry_InternalError(t *testing.T) {
	callNumber := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callNumber++
		http.Error(w, "test error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := GoogleTranslateClient{
		endpoint: ts.URL,
		apiKey:   "test",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		client.endpoint,
		nil,
	)
	require.NoError(t, err)
	responseBody, err := client.GetResponseWithRetry(request)
	require.Nil(t, responseBody)
	require.NotNil(t, err)
	require.Equal(t, 2, callNumber)
}
