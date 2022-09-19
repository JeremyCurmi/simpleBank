package database

import (
	"errors"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/jmoiron/sqlx"
)

type DBUser struct {
	db *sqlx.DB
}

func NewDBUser(db *sqlx.DB) *DBUser {
	return &DBUser{
		db: db,
	}
}

func (s *DBUser) Create(model models.DBModel) error {
	query := "INSERT INTO users (username, password) VALUES (:username, :password)"
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
func (s *DBUser) Update(queryParams ...interface{}) error {
	return nil
}
func (s *DBUser) Delete(queryParams ...interface{}) error {
	return nil
}
func (s *DBUser) Get(queryParams ...interface{}) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password, created_at FROM users WHERE id = $1`
	if err := s.db.Get(&user, query, queryParams...); err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (s *DBUser) GetByUsername(username string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password, created_at FROM users WHERE username = $1`
	if err := s.db.Get(&user, query, username); err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (s *DBUser) GetAll(queryParams ...interface{}) ([]models.User, error) {
	return nil, nil
}
