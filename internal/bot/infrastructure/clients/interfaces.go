package clients

import (
	"context"
)

type TranslationClient interface {
	GetTranslation(ctx context.Context, source, target, text string) (string, error)
}

type CoordinateConverter interface {
	ConvertETRSGK24toWGS84(
		ctx context.Context,
		latitude float32,
		longitude float32,
	) (float64, float64, error)
}
