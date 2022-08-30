package utils

import (
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
)

func RunMigrations(logger *zap.Logger, migrationURL, dbSource string) error {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("No changes detected for migrations")
			return nil
		}
		return err
	}
	logger.Info("Migrations ran successfully ✅")
	return nil
}
