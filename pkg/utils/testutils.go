package utils

import (
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"gitlab.com/go_projects_jer/simple_bank/pkg/config"
)

func SetUpTestDB(t *testing.T) *sqlx.DB {
	logger := NewLogger()
	testDBClient, err := sqlx.Connect("postgres", *config.DBURL)
	require.NoError(t, err)
	require.NoError(t, RunMigrations(logger, "file://../../pkg/migrations", *config.DBURL))
	return testDBClient
}
