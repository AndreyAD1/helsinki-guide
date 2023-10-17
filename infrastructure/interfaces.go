package infrastructure

import "context"


type translationClient interface {
	GetTranslation(context.Context, string, string, string) (string, error)
}