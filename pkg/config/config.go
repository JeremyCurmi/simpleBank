package config

import (
	"flag"

	"github.com/alex-ant/envs"
)

var (
	LogLevel = flag.String("log-level", "debug", "log level")
	DBURL = flag.String("db-url", "postgres://postgres:demo_password@localhost:5432/Bank?sslmode=disable", "database connection url")
	MigrationURL = flag.String("migration-url", "file://pkg/migrations", "path to migrations")
	APIPort = flag.String("api-port", "3000", "api port")
	APIHost = flag.String("api-host", "0.0.0.0", "api host")
	DBConnMaxLifetime = flag.Int("db-conn-max-lifetime", 60, "database connection max lifetime")
	TokenHourLifeSpan = flag.Int("token-hour-lifespan", 10, "token hour life span")
	APISecretKey = flag.String("api-secret-key", "secret", "api secret key")
)

func InitEnvVariables() error {
	flag.Parse()
	if err := envs.GetAllFlags(); err!=nil {
		return err
	}
	return nil
}