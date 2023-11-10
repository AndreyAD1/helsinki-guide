package integrationtests

import (
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"

	"github.com/stretchr/testify/require"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestMigrations(t *testing.T) {
	migrationPath := filepath.Join(
		"..",
		"internal",
		"infrastructure",
		"migrations",
	)
	m, err := migrate.New("file:"+migrationPath, databaseUrl)
	require.NoErrorf(
		t,
		err,
		"can not instantiate a migration tool '%s' for '%s': %v",
		migrationPath,
		databaseUrl,
		err,
	)
	require.NoError(t, m.Up())
	require.NoError(t, m.Down())
	require.NoError(t, m.Up())
}
