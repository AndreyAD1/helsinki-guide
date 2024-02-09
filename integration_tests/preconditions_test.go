package integrationtests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
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
	dbpool         *pgxpool.Pool
	databaseUrl    string
	logLevel       = new(slog.LevelVar)
	MIGRATION_PATH = filepath.Join(
		"..",
		"internal",
		"bot",
		"infrastructure",
		"migrations",
	)
)

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION") == "" {
		log.Println("SKIP integration tests: set an 'INTEGRATION' environment variable")
		return
	}
	log.Println("run integration tests")
	handlerOptions := slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &handlerOptions))
	slog.SetDefault(logger)
	logLevel.Set(slog.LevelDebug)

	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct a pool: %s", err)
	}

	log.Println("ping a Docker service")
	err = dockerPool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	log.Println("run a PostgreSQL container")
	resource, err := dockerPool.BuildAndRunWithBuildOptions(
		&dockertest.BuildOptions{
			Dockerfile: "./Dockerfile",
			ContextDir: ".",
			BuildArgs:  []docker.BuildArg{{Name: "tag", Value: "guide_test_db"}},
		},
		&dockertest.RunOptions{
			Name: "guide_test_db",
			Env: []string{
				"POSTGRES_PASSWORD=secret",
				"POSTGRES_USER=user_name",
				"POSTGRES_DB=dbname",
				"listen_addresses = '*'",
			},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		log.Fatalf("Could not start a resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl = fmt.Sprintf(
		"postgres://user_name:secret@%s/dbname?sslmode=disable",
		hostAndPort,
	)

	log.Println("Connecting to the database on url: ", databaseUrl)

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
		log.Fatalf("Could not connect to a DB container: %s", err)
	}
	var code int
	defer func() {
		log.Printf("a test exit code: %v; cleaning up", code)
		if err := dockerPool.Purge(resource); err != nil {
			log.Fatalf("Could not purge a docker resource: %s", err)
		}
		os.Exit(code)
	}()
	code = m.Run()
}

func TestDBInteractions(t *testing.T) {
	log.Println("run DB tests")
	migrator, err := migrate.New("file:"+MIGRATION_PATH, databaseUrl)
	require.NoErrorf(
		t,
		err,
		"can not instantiate a migration tool '%s' for '%s': %v",
		MIGRATION_PATH,
		databaseUrl,
		err,
	)
	t.Cleanup(func() { migrator.Drop() })
	for _, test := range integrationTests {
		err = migrator.Up()
		errCheck := func() bool {
			if err == nil || errors.Is(err, migrate.ErrNoChange) {
				return true
			}
			return false
		}
		require.Conditionf(
			t,
			errCheck,
			fmt.Sprintf("a migration error for the test '%v': %v", test.name, err),
		)
		t.Run(test.name, test.function)
		migrator.Down()
	}
}

type integrationTest struct {
	name     string
	function func(*testing.T)
}

var integrationTests = []integrationTest{
	{"neigbourhoods", testNeighbourhoodRepository},
	{"actors", testActorRepository},
	{"addBuilding", testAddNewBuilding},
	{"addBuildingAddressError", testAddNewBuildingAddressError},
	{"addBuildingAuthorError", testAddNewBuildingAuthorError},
	{"getNearestBuildings", testGetNearestBuildings},
	{"updateAbsentBuilding", testUpdateAbsentBuilding},
	{"manageRemovedBuilding", testManageRemovedBuilding},
	{"runPopulator", testRunPopulator},
	{"addUser", testUserRepository},
}
