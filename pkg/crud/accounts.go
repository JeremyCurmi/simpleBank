package crud

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AccountsService struct {
	logger *zap.Logger
	db *sqlx.DB
}

func NewAccountsService(logger *zap.Logger, db *sqlx.DB) *AccountsService {
	return &AccountsService{logger, db}
}