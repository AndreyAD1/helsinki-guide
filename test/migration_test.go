package integrationtests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	dbpool      *pgxpool.Pool
	databaseUrl string
)

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION") == "" {
		log.Println("skip integration tests: set INTEGRATION environment variable")
		return
	}

	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct a pool: %s", err)
	}

	err = dockerPool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start a resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl = fmt.Sprintf(
		"postgres://user_name:secret@%s/dbname?sslmode=disable",
		hostAndPort,
	)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120)

	dockerPool.MaxWait = 120 * time.Second
	ctx := context.Background()
	if err = dockerPool.Retry(func() error {
		dbpool, err = pgxpool.New(context.Background(), databaseUrl)
		if err != nil {
			log.Printf(
				"unable to create a connection pool: DB URL '%v': %v",
				databaseUrl,
				err,
			)
			os.Exit(1)
		}
		log.Println("ping a DB")
		return dbpool.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	if err := dockerPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge a resource: %s", err)
	}

	os.Exit(code)
}

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
		"can not start migrations '%s' for '%s': %v",
		migrationPath,
		databaseUrl,
		err,
	)
	require.NoError(t, m.Up())
	require.NoError(t, m.Down())
	require.NoError(t, m.Up())
}
