package crud

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"gitlab.com/go_projects_jer/simple_bank/pkg/models"
	"gitlab.com/go_projects_jer/simple_bank/pkg/utils"
	"go.uber.org/zap"
)

type AuthService struct {
	logger *zap.Logger
	db *sqlx.DB
}

func NewAuthService(logger *zap.Logger, db *sqlx.DB) *AuthService {
	return &AuthService{logger, db}
}

func (s *AuthService) ValidateUser(user *models.UserAPI) (string, error) {

	userDB, err := s.GetUser(user.UserName)
	if err != nil {
		return "", err
	}
	if !utils.CheckPasswordHash(user.Password, userDB.Password) {
		return "", errors.New(utils.ErrIncorrectPassword)
	}

	token, err := utils.GenerateToken(userDB.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) CreateUser(user *models.UserAPI) error {
	var err error
	stmt := "INSERT INTO users (username, password) VALUES (:username, :password)"
	_, err = s.db.NamedExec(stmt, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetUser(username string) (*models.User, error) {
	var user models.User
	err := s.db.Get(&user, `SELECT id, username, password, created_at FROM users WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) CurrentUser(c *gin.Context) (*models.User, error) {

	userID, err := utils.ExtractTokenID(c)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}