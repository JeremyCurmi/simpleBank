package api

import (
	"fmt"
	"github.com/JeremyCurmi/simpleBank/pkg/database"
	"github.com/JeremyCurmi/simpleBank/pkg/services"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	user     = "postgres"
	password = "secret"
	db       = "postgres"
	//port              = "5433"
	dsn               = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	dbConn            *sqlx.DB
	dbConnMaxLifetime = 60
	logger            = utils.NewLogger()
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
		//ExposedPorts: []string{"5432"},
		//PortBindings: map[docker.Port][]docker.PortBinding{
		//	"5432": {
		//		{HostIP: "0.0.0.0", HostPort: port},
		//	},
		//},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	dsn = fmt.Sprintf(dsn, user, password, resource.GetPort("5432/tcp"), db)
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

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatalf("could not close DB connection: %s", err)
		}
	}()

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func setup() *Manager {
	logger := utils.NewLogger()
	dbClient, err := database.New(logger, dsn, dbConnMaxLifetime)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}
	dbConn := dbClient.Conn()
	userService := services.NewUserService(logger, database.NewDBUser(dbConn))
	accountsService := services.NewAccountsService(logger, database.NewDBAccount(dbConn))
	transferService := services.NewTransferService(logger, database.NewDBTransfer(dbConn), accountsService)
	return New(logger, accountsService, userService, transferService)
}

func SetupRouter() *gin.Engine {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	return r
}
