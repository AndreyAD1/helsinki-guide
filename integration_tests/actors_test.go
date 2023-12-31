package integrationtests

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/stretchr/testify/require"
)

func testActorRepository(t *testing.T) {
	storage := repositories.NewActorRepo(dbpool)
	titleEn := "test title en"
	actor := repositories.Actor{Name: "test", TitleEn: &titleEn}
	saved, err := storage.Add(context.Background(), actor)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)

	spec := repositories.NewActorSpecificationByName(actor)
	stored, err := storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equal(t, 1, len(stored))
	require.Equal(t, *saved, stored[0])

	saved2, err := storage.Add(context.Background(), actor)
	require.ErrorIs(t, err, repositories.ErrDuplicate)
	require.Equal(t, saved, saved2)
}
