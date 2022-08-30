package main

import (
	"log"

	"gitlab.com/go_projects_jer/simple_bank/pkg/api"
	"gitlab.com/go_projects_jer/simple_bank/pkg/config"
	"gitlab.com/go_projects_jer/simple_bank/pkg/database"
	"gitlab.com/go_projects_jer/simple_bank/pkg/utils"
	"go.uber.org/zap"
)

func init() {
	err := config.InitEnvVariables()
	if err != nil {
		log.Fatalf("Error getting env variables: %s", err)
	}
}

func main() {
	logger := utils.NewLogger()
	dbClient, err := database.New(logger, *config.DBURL, *config.DBConnMaxLifetime)
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = utils.RunMigrations(logger, *config.MigrationURL, *config.DBURL)
	if err != nil {
		logger.Fatal(utils.ErrDBMigrations, zap.Error(err))
	}

	err = api.InitAPI(logger, dbClient.Conn())
	if err != nil {
		logger.Fatal(utils.ErrApiInitial, zap.Error(err))
	}
}
