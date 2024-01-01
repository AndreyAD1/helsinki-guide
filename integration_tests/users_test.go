package integrationtests

import (
	"context"
	"testing"

	r "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func testUserRepository(t *testing.T) {
	storage := r.NewUserRepo(dbpool)
	user := r.User{TelegramID: 123, PreferredLanguage: "fi"}
	saved, err := storage.AddOrUpdate(context.Background(), user)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)

	spec := r.NewUserSpecificationByID(123)
	stored, err := storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equal(t, 1, len(stored))
	require.Equal(
		t,
		cmp.Diff(
			*saved,
			stored[0],
			cmpopts.IgnoreUnexported(r.Timestamps{}),
		),
		"",
	)

	saved.PreferredLanguage = "en"
	updated, err := storage.AddOrUpdate(context.Background(), *saved)
	require.NoError(t, err)
	require.Equal(
		t,
		cmp.Diff(
			saved,
			updated,
			cmpopts.IgnoreUnexported(r.Timestamps{}),
			cmpopts.IgnoreFields(r.User{}, "UpdatedAt"),
		),
		"",
	)

	stored2, err := storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equal(t, 1, len(stored2))
	require.Equal(t, *updated, stored2[0])
}
