package services

import (
	"context"
	"fmt"
	"github.com/JeremyCurmi/simpleBank/pkg/database"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"go.uber.org/zap"
)

type TransferService struct {
	logger          *zap.Logger
	dbService       *database.DBTransfer
	accountsService *AccountsService
}

var txKey = struct{}{}

func NewTransferService(logger *zap.Logger, dbService *database.DBTransfer, accountsService *AccountsService) *TransferService {
	return &TransferService{
		logger:          logger,
		dbService:       dbService,
		accountsService: accountsService,
	}
}
func (s *TransferService) CreateTransfer(transfer models.Transfer) error {
	if err := s.dbService.Create(transfer); err != nil {
		s.logger.Error("error creating transfer", zap.Error(err))
	}
	return nil
}
func (s *TransferService) GetTransfersBySender(id uint) ([]models.Transfer, error) {
	transfers, err := s.dbService.GetBySender(id)
	if err != nil {
		s.logger.Error("error getting transfers by sender", zap.Error(err))
		return []models.Transfer{}, err
	}
	return transfers, nil
}
func (s *TransferService) GetTransfersByReceiver(id uint) ([]models.Transfer, error) {
	transfers, err := s.dbService.GetByReceiver(id)
	if err != nil {
		s.logger.Error("error getting transfers by receiver", zap.Error(err))
		return []models.Transfer{}, err
	}
	return transfers, nil
}
func (s *TransferService) GetTransfersBySenderAndReceiver(senderID, ReceiverID uint) ([]models.Transfer, error) {
	transfers, err := s.dbService.GetBySenderAndReceiver(senderID, ReceiverID)
	if err != nil {
		s.logger.Error("error getting transfers by sender and receiver", zap.Error(err))
		return []models.Transfer{}, err
	}
	return transfers, nil
}
func (s *TransferService) TransferFunds(ctx context.Context, transfer *models.Transfer) (error, error) {
	fail := func(err error) error {
		return fmt.Errorf("transfer of funds failed: %v", err)
	}
	transaction := func() error {
		tx, err := s.dbService.BeginTx(ctx, nil)
		if err != nil {
			s.logger.Error("failed to begin transaction")
			return fail(err)
		}
		txName := ctx.Value(txKey)
		defer func() {
			if err := tx.Rollback(); err != nil {
				s.logger.Fatal(fmt.Sprintf("failed to rollback transaction: %v", err))
			}
		}()

		if transfer.Amount <= 0 {
			return fail(fmt.Errorf("amount must be greater than 0"))
		}
		senderAccount, err := s.accountsService.GetAccount(transfer.SenderID, true, ctx)
		if err != nil {
			s.logger.Error("failed to get sender account")
			return fail(err)
		}
		fmt.Println(txName, "get sender", senderAccount.Balance)

		// TODO: Create transaction for sender
		// TODO: Create transaction for receiver

		if senderAccount.Balance < transfer.Amount {
			s.logger.Error("sender does not have enough balance")
			return fail(fmt.Errorf("sender does not have enough balance"))
		}
		senderAccount.Balance -= transfer.Amount
		fmt.Println(txName, "update sender")
		err = s.accountsService.UpdateAccountBalance(&senderAccount, transfer.SenderID)
		if err != nil {
			s.logger.Error("failed to update sender account")
			return fail(err)
		}

		fmt.Println(txName, "get receiver")
		receiverAccount, err := s.accountsService.GetAccount(transfer.ReceiverID, true, ctx)
		if err != nil {
			s.logger.Error("failed to get receiver account")
			return fail(err)
		}
		receiverAccount.Balance += transfer.Amount
		fmt.Println(txName, "update receiver")
		err = s.accountsService.UpdateAccountBalance(&receiverAccount, transfer.ReceiverID)
		if err != nil {
			s.logger.Error("failed to update receiver account")
			return fail(err)
		}

		fmt.Println(txName, "create transfer")
		transfer.Status = "success"
		if err := s.CreateTransfer(*transfer); err != nil {
			s.logger.Error("failed to create transfer record")
			return fail(err)
		}
		return nil
	}

	if err := transaction(); err != nil {
		transfer.Status = "failed"
		return err, s.CreateTransfer(*transfer)
	}
	return nil, nil
}
