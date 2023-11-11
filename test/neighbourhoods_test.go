package integrationtests

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
	"github.com/stretchr/testify/require"
)

func testNeighbourhoodRepository(t *testing.T) {
	storage := repositories.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := internal.Neighbourhood{Name: "test"}
	saved, err := storage.Add(context.Background(), neighbourbourhood)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)
}
