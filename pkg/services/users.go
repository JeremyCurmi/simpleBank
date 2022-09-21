package services

import (
	"errors"
	"fmt"
	"github.com/JeremyCurmi/simpleBank/pkg/database"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserService struct {
	logger    *zap.Logger
	dbService *database.DBUser
}

func NewUserService(logger *zap.Logger, dbService *database.DBUser) *UserService {
	return &UserService{
		logger:    logger,
		dbService: dbService,
	}
}
func (s *UserService) ValidateUser(user *models.User) (string, error) {
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
func (s *UserService) CreateUser(user *models.User) error {
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("error hashing password: %v", zap.Error(err))
		return err
	}

	user.Password = string(hashPassword)
	err = s.dbService.Create(user)
	if err != nil {
		s.logger.Error("error creating user: %v", zap.Error(err))
		return err
	}
	return nil
}
func (s *UserService) GetUser(username string) (models.User, error) {
	user, err := s.dbService.GetByUsername(username)
	if err != nil {
		s.logger.Error("error getting user: %v", zap.Error(err))
		return models.User{}, err
	}
	return user, nil
}
func (s *UserService) GetUserByID(id uint) (models.User, error) {
	fmt.Println("user by id:", id)
	user, err := s.dbService.Get(id)
	if err != nil {
		s.logger.Error("error getting user by id: %v", zap.Error(err))
		return models.User{}, err
	}
	return user, err
}
func (s *UserService) CurrentUser(c *gin.Context) (*models.User, error) {
	userID, err := utils.ExtractTokenID(c)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
