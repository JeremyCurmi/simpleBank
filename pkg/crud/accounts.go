package crud

import (
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AccountsService struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func NewAccountsService(logger *zap.Logger, db *sqlx.DB) *AccountsService {
	return &AccountsService{logger, db}
}

func (s *AccountsService) CreateAccount(account models.Account) (uint, error) {
	stmt := `INSERT INTO accounts (user_id, name, balance, currency) VALUES (:user_id, :name, :balance, :currency)`
	result, err := s.db.NamedExec(stmt, account)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

func (s *AccountsService) UpdateUserAccountBalance(id uint, Account models.Account) (int, error) {
	stmt := `UPDATE accounts SET balance = $1  WHERE id = $2`
	result, err := s.db.Exec(stmt, Account.Balance, id)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *AccountsService) UpdateUserAccountName(id uint, Account models.Account) (int, error) {
	stmt := `UPDATE accounts SET name = $1  WHERE id = $2`
	result, err := s.db.Exec(stmt, Account.Name, id)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *AccountsService) GetUserAccounts(UserID uint) ([]models.Account, error) {
	var accounts []models.Account
	stmt := `SELECT id, name, user_id, balance, currency, created_at, updated_at FROM accounts WHERE user_id = $1`
	err := s.db.Select(&accounts, stmt, UserID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *AccountsService) GetAccount(UserID uint, Name string) (models.Account, error) {
	var account models.Account
	stmt := `SELECT id, name, user_id, balance, currency, created_at, updated_at FROM accounts WHERE user_id = $1 AND name = $2`
	err := s.db.Get(&account, stmt, UserID, Name)
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func (s *AccountsService) GetAccountByID(id uint) (models.Account, error) {
	var account models.Account
	stmt := `SELECT id, name, user_id, balance, currency, created_at, updated_at FROM accounts WHERE id = $1`
	err := s.db.Get(&account, stmt, id)
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func (s *AccountsService) DeleteAccount(UserID uint, name string) (int, error) {
	stmt := `DELETE FROM accounts WHERE user_id = $1 AND name = $2`
	res, err := s.db.Exec(stmt, UserID, name)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), err
}
