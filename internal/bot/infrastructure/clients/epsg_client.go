package clients

import "context"

type EPSGClient struct {
	baseURL string
}

func NewEPSGClient(baseURL string) *EPSGClient {
	return &EPSGClient{baseURL}
}

func (c *EPSGClient) ConvertETRSGK24toWGS84(
	ctx context.Context,
	latitude,
	longitude float32,
) (float32, float32, error) {
	return 0, 0, nil
}
