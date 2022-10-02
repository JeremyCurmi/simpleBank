package services

import (
	"context"
	"errors"
	"github.com/JeremyCurmi/simpleBank/pkg/database"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"go.uber.org/zap"
)

type AccountsService struct {
	logger    *zap.Logger
	dbService *database.DBAccount
}

func NewAccountsService(logger *zap.Logger, dbService *database.DBAccount) *AccountsService {
	return &AccountsService{
		logger:    logger,
		dbService: dbService,
	}
}
func (s *AccountsService) CreateAccount(account models.Account) error {
	if err := s.dbService.Create(account); err != nil {
		s.logger.Error("error creating account", zap.Error(err))
	}
	return nil
}
func (s *AccountsService) UpdateAccountBalance(model *models.Account, id uint) error {
	if model.Balance <= 0 {
		return errors.New("balance must be positive")
	}
	if err := s.dbService.Update(model.Balance, id); err != nil {
		s.logger.Error("error updating account", zap.Error(err))
		return err
	}
	return nil
}
func (s *AccountsService) UpdateAccountName(model *models.Account, id uint) error {
	if model.Name == "" {
		return errors.New("account name is required")
	}

	if err := s.dbService.UpdateName(model.Name, id); err != nil {
		s.logger.Error("error updating account", zap.Error(err))
		return err
	}
	return nil
}
func (s *AccountsService) GetUserAccounts(userID uint) ([]models.Account, error) {
	accounts, err := s.dbService.GetAll(userID)
	if err != nil {
		s.logger.Error("error getting user accounts", zap.Error(err))
		return []models.Account{}, err
	}
	return accounts, nil
}
func (s *AccountsService) GetAccount(id uint, lock bool, ctx context.Context) (models.Account, error) {
	var (
		account models.Account
		err     error
	)
	if lock {
		account, err = s.dbService.GetLock(ctx, id)
	} else {
		account, err = s.dbService.Get(id)
	}

	if err != nil {
		s.logger.Error("error getting account", zap.Error(err))
		return models.Account{}, err
	}
	return account, nil
}
func (s *AccountsService) GetAccountByName(userID uint, name string) (models.Account, error) {
	account, err := s.dbService.GetByUserIDAndName(userID, name)
	if err != nil {
		s.logger.Error("error getting account", zap.Error(err))
		return models.Account{}, err
	}
	return account, err

}
func (s *AccountsService) DeleteAccount(id uint) error {
	if err := s.dbService.Delete(id); err != nil {
		s.logger.Error("error deleting account", zap.Error(err))
		return err
	}
	return nil
}
