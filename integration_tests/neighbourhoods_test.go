package integrationtests

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	s "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	"github.com/stretchr/testify/require"
)

func testNeighbourhoodRepository(t *testing.T) {
	storage := repositories.NewNeighbourhoodRepo(dbpool)
	neighbourhood := types.Neighbourhood{Name: "test"}
	saved, err := storage.Add(context.Background(), neighbourhood)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)

	spec := s.NewNeighbourhoodSpecificationByName(neighbourhood)
	stored, err := storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equal(t, 1, len(stored))
	require.Equal(t, *saved, stored[0])

	saved2, err := storage.Add(context.Background(), neighbourhood)
	require.ErrorIs(t, err, repositories.ErrDuplicate)
	require.Equal(t, saved, saved2)

	municipality := "Helsinki"
	neighbourhood = types.Neighbourhood{
		Name:         "test",
		Municipality: &municipality,
	}
	saved, err = storage.Add(context.Background(), neighbourhood)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)

	spec = s.NewNeighbourhoodSpecificationAll(100, 0)
	stored, err = storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equal(t, 2, len(stored))
}
