package clients

import (
	"context"
)

type TranslationClient interface {
	GetTranslation(ctx context.Context, source, target, text string) (string, error)
}
