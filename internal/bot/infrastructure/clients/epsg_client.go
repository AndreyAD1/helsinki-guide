package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	EPSGcodeETRSGK25 = 3879
	EPSGcodeWGS84    = 4326
)

type EPSGClient struct {
	baseURL string
	timeout int
	client  *http.Client
}

type PointResponse struct {
	Latitude  string `json:"y"`
	Longitude string `json:"x"`
	Height    string `json:"z"`
}

func NewEPSGClient(baseURL string, timeoutSeconds int) *EPSGClient {
	return &EPSGClient{baseURL, timeoutSeconds, &http.Client{}}
}

func (c *EPSGClient) ConvertETRSGK24toWGS84(
	ctx context.Context,
	latitude,
	longitude float32,
) (float64, float64, error) {
	requestURL, err := url.JoinPath(c.baseURL, "trans")
	if err != nil {
		return 0, 0, err
	}
	query := fmt.Sprintf(
		"x=%.2f&y=%.2f&s_srs=%v&t_srs=%v",
		longitude,
		latitude,
		EPSGcodeETRSGK25,
		EPSGcodeWGS84,
	)
	ctx, cancel := context.WithTimeout(
		ctx,
		time.Second*time.Duration(c.timeout),
	)
	defer cancel()
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL+"?"+query,
		bytes.NewBuffer([]byte{}),
	)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"can not make a request for '%.2f': '%.2f': %v",
			latitude,
			longitude,
			err,
		)
	}
	response, err := GetResponseWithRetry(c.client, request)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"can not send a request for '%.2f': '%.2f': %v",
			latitude,
			longitude,
			err,
		)
	}
	var point PointResponse
	if err = json.Unmarshal(response, &point); err != nil {
		err = fmt.Errorf(
			"receive an unexpected response body: '%v': %w",
			string(response),
			err,
		)
		return 0, 0, err
	}
	convertedLatitude, err := strconv.ParseFloat(point.Latitude, 64)
	if err != nil {
		return 0, 0, err
	}
	convertedLongitude, err := strconv.ParseFloat(point.Longitude, 64)
	if err != nil {
		return 0, 0, err
	}
	return convertedLatitude, convertedLongitude, nil
}
