package crud

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/go_projects_jer/simple_bank/pkg/models"
	"go.uber.org/zap"
)

type AccountsService struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func NewAccountsService(logger *zap.Logger, db *sqlx.DB) *AccountsService {
	return &AccountsService{logger, db}
}

func (s *AccountsService) CreateAccount(account models.Account) (int, error) {
	stmt := `INSERT INTO accounts (user_id, name, balance, currency) VALUES (:user_id, :name, :balance, :currency)`
	result, err := s.db.NamedExec(stmt, account)
	if err != nil {
		return 0, err
	}

	accountID, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(accountID), nil
}

func (s *AccountsService) UpdateUserAccount(UserID uint, name string, account models.Account) (int, error) {
	return 0, nil
}

func (s *AccountsService) GetUserAccounts(UserID uint) ([]models.Account, error) {
	return []models.Account{}, nil
}

func (s *AccountsService) GetAccount(UserID uint, name string) (models.Account, error) {
	return models.Account{}, nil
}

func (s *AccountsService) DeleteAccount(UserID uint, name string) error {
	return nil
}
