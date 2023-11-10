package integrationtests

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
	"github.com/stretchr/testify/require"
)

func TestActorRepository(t *testing.T) {
	storage := repositories.NewActorRepo(dbpool)
	titleEn := "test title en"
	actor := internal.Actor{Name: "test", TitleEn: &titleEn}
	saved, err := storage.Add(context.Background(), actor)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)
}
