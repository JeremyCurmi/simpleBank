package database

import (
	"context"
	"errors"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/jmoiron/sqlx"
)

type DBAccount struct {
	db *sqlx.DB
}

func NewDBAccount(db *sqlx.DB) *DBAccount {
	return &DBAccount{
		db: db,
	}
}
func (s *DBAccount) Create(model models.DBModel) error {
	query := `INSERT INTO accounts (user_id, name, balance, currency) VALUES (:user_id, :name, :balance, :currency)`
	result, err := s.db.NamedExec(query, model)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no rows were inserted")
	}

	return nil
}
func (s *DBAccount) Update(queryParams ...interface{}) error {
	query := `UPDATE accounts SET balance = $1 WHERE id = $2`
	result, err := s.db.Exec(query, queryParams...)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no rows were updated")
	}

	return nil
}
func (s *DBAccount) UpdateName(queryParams ...interface{}) error {
	query := `UPDATE accounts SET name = $1 WHERE id = $2`
	result, err := s.db.Exec(query, queryParams...)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no rows were updated")
	}

	return nil
}
func (s *DBAccount) Delete(queryParams ...interface{}) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := s.db.Exec(query, queryParams...)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no rows were deleted")
	}

	return nil
}
func (s *DBAccount) Get(queryParams ...interface{}) (models.Account, error) {
	var account models.Account
	query := `SELECT id, name, user_id, balance, currency, created_at, updated_at FROM accounts WHERE id = $1`
	if err := s.db.Get(&account, query, queryParams...); err != nil {
		return models.Account{}, err
	}
	return account, nil
}
func (s *DBAccount) GetLock(ctx context.Context, queryParams ...interface{}) (models.Account, error) {
	var account models.Account
	query := `SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR UPDATE;`
	if err := s.db.GetContext(ctx, &account, query, queryParams...); err != nil {
		return models.Account{}, err
	}
	return account, nil
}
func (s *DBAccount) GetByUserIDAndName(queryParams ...interface{}) (models.Account, error) {
	var account models.Account
	query := `SELECT id, name, user_id, balance, currency, created_at, updated_at FROM accounts WHERE user_id = $1 and name = $2`
	if err := s.db.Get(&account, query, queryParams...); err != nil {
		return models.Account{}, err
	}
	return account, nil
}
func (s *DBAccount) GetAll(queryParams ...interface{}) ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT id, name, user_id, balance, currency, created_at, updated_at FROM accounts WHERE user_id = $1`
	if err := s.db.Select(&accounts, query, queryParams...); err != nil {
		return nil, err
	}
	return accounts, nil
}
