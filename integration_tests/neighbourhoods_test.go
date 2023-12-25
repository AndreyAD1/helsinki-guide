package integrationtests

import (
	"context"
	"testing"

	r "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/stretchr/testify/require"
)

func testNeighbourhoodRepository(t *testing.T) {
	storage := r.NewNeighbourhoodRepo(dbpool)
	neighbourhood := r.Neighbourhood{Name: "test"}
	saved, err := storage.Add(context.Background(), neighbourhood)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)

	spec := r.NewNeighbourhoodSpecificationByName(neighbourhood)
	stored, err := storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equal(t, 1, len(stored))
	require.Equal(t, *saved, stored[0])

	saved2, err := storage.Add(context.Background(), neighbourhood)
	require.ErrorIs(t, err, r.ErrDuplicate)
	require.Equal(t, saved, saved2)

	municipality := "Helsinki"
	neighbourhood = r.Neighbourhood{
		Name:         "test",
		Municipality: &municipality,
	}
	saved, err = storage.Add(context.Background(), neighbourhood)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)

	spec = r.NewNeighbourhoodSpecificationAll(100, 0)
	stored, err = storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equal(t, 2, len(stored))
}
