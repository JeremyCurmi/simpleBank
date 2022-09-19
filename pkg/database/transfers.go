package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/jmoiron/sqlx"
)

type DBTransfer struct {
	db *sqlx.DB
}

func NewDBTransfer(db *sqlx.DB) *DBTransfer {
	return &DBTransfer{db: db}
}
func (s *DBTransfer) Create(model models.DBModel) error {
	query := `INSERT INTO transfers (sender_id, receiver_id, amount, status) VALUES (:sender_id, :receiver_id, :amount, :status)`
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
func (s *DBTransfer) Update(queryParams ...interface{}) error {
	return nil
}
func (s *DBTransfer) Delete(queryParams ...interface{}) error { return nil }
func (s *DBTransfer) Get(queryParams ...interface{}) (models.Transfer, error) {
	var transfer models.Transfer
	query := `SELECT * FROM transfers WHERE id = $1`
	if err := s.db.Get(&transfer, query, queryParams...); err != nil {
		return models.Transfer{}, err
	}
	return transfer, nil
}
func (s *DBTransfer) GetBySender(queryParams ...interface{}) ([]models.Transfer, error) {
	var transfer []models.Transfer
	query := `SELECT * FROM transfers WHERE sender_id = $1`
	if err := s.db.Select(&transfer, query, queryParams...); err != nil {
		return []models.Transfer{}, err
	}
	return transfer, nil
}
func (s *DBTransfer) GetByReceiver(queryParams ...interface{}) ([]models.Transfer, error) {
	var transfer []models.Transfer
	query := `SELECT * FROM transfers WHERE receiver_id = $1`
	if err := s.db.Select(&transfer, query, queryParams...); err != nil {
		return []models.Transfer{}, err
	}
	return transfer, nil
}
func (s *DBTransfer) GetBySenderAndReceiver(senderID, receiverID uint) ([]models.Transfer, error) {
	var transfer []models.Transfer
	query := `SELECT * FROM transfers WHERE sender_id = $1 and receiver_id = $2`
	if err := s.db.Select(&transfer, query, senderID, receiverID); err != nil {
		return []models.Transfer{}, err
	}
	return transfer, nil
}
func (s *DBTransfer) BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, options)
}
