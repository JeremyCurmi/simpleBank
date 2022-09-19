package services

import (
	"context"
	"github.com/JeremyCurmi/simpleBank/pkg/database"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func setupTransfers() (*AccountsService, *TransferService) {
	dbAccount := database.NewDBAccount(dbConn)
	dbTransfer := database.NewDBTransfer(dbConn)
	accountService := NewAccountsService(logger, dbAccount)
	transferService := NewTransferService(logger, dbTransfer, accountService)
	return accountService, transferService
}

func TestTransferFunds(t *testing.T) {
	accountService, transferService := setupTransfers()
	account1Name := utils.RandomAccountName()
	account2Name := utils.RandomAccountName()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(200*time.Second))
	defer cancel()
	accounts := []models.Account{
		{
			UserID:   testUserID,
			Name:     account1Name,
			Currency: userCurrency,
			Balance:  20,
		},
		{
			UserID:   adminUserID,
			Name:     account2Name,
			Currency: userCurrency,
			Balance:  10,
		},
	}

	for _, account := range accounts {
		err := accountService.CreateAccount(account)
		require.NoError(t, err)
	}

	senderAccount, err := accountService.GetAccountByName(testUserID, account1Name)
	require.NoError(t, err)

	receiverAccount, err := accountService.GetAccountByName(adminUserID, account2Name)
	require.NoError(t, err)

	t.Run("sender has enough funds", func(t *testing.T) {
		transferData := models.Transfer{
			SenderID:   senderAccount.ID,
			ReceiverID: receiverAccount.ID,
			Amount:     10,
		}
		err := transferService.TransferFunds(ctx, &transferData)
		require.NoError(t, err)

		senderAccount, _ = accountService.GetAccountByName(testUserID, account1Name)
		require.Equal(t, float64(10), senderAccount.Balance)

		receiverAccount, _ = accountService.GetAccountByName(adminUserID, account2Name)
		require.Equal(t, float64(20), receiverAccount.Balance)
	})
	t.Run("sender does not have enough funds", func(t *testing.T) {
		transferData := models.Transfer{
			SenderID:   senderAccount.ID,
			ReceiverID: receiverAccount.ID,
			Amount:     100,
		}
		err := transferService.TransferFunds(ctx, &transferData)
		require.Error(t, err)
		require.Equal(t, "transfer of funds failed: sender does not have enough balance", err.Error())
	})
	t.Run("sender sends negative amount", func(t *testing.T) {
		transferData := models.Transfer{
			SenderID:   testUserID,
			ReceiverID: adminUserID,
			Amount:     -10,
		}
		err := transferService.TransferFunds(ctx, &transferData)
		require.Error(t, err)
		require.Equal(t, "transfer of funds failed: amount must be greater than 0", err.Error())
	})
}