package crud

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"gitlab.com/go_projects_jer/simple_bank/pkg/utils"
)

var (
	user     = "postgres"
	password = "secret"
	db       = "postgres"
	port     = "5433"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	dbConn   *sqlx.DB
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	dsn = fmt.Sprintf(dsn, user, password, port, db)
	if err = pool.Retry(func() error {
		dbConn, err = sqlx.Connect("postgres", dsn)
		if err != nil {
			return err
		}
		err = utils.RunMigrations(logger, "file://../../pkg/migrations", dsn)
		return err

	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	defer dbConn.Close()

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
