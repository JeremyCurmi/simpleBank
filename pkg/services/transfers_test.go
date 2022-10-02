package services

import (
	"context"
	"fmt"
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
		txErr, insertErr := transferService.TransferFunds(ctx, &transferData)
		require.NoError(t, txErr)
		require.NoError(t, insertErr)

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
		txErr, insertErr := transferService.TransferFunds(ctx, &transferData)
		require.NoError(t, insertErr)
		require.Error(t, txErr)
		require.Equal(t, "transfer of funds failed: sender does not have enough balance", txErr.Error())
	})
	t.Run("sender sends negative amount", func(t *testing.T) {
		transferData := models.Transfer{
			SenderID:   testUserID,
			ReceiverID: adminUserID,
			Amount:     -10,
		}
		txErr, insertErr := transferService.TransferFunds(ctx, &transferData)
		require.NoError(t, insertErr)
		require.Error(t, txErr)
		require.Equal(t, "transfer of funds failed: amount must be greater than 0", txErr.Error())
	})
}

func TestConcurrentTransfers(t *testing.T) {
	accountService, transferService := setupTransfers()
	account1Name := utils.RandomAccountName()
	account2Name := utils.RandomAccountName()
	accounts := []models.Account{
		{
			UserID:   testUserID,
			Name:     account1Name,
			Currency: userCurrency,
			Balance:  100,
		},
		{
			UserID:   adminUserID,
			Name:     account2Name,
			Currency: userCurrency,
			Balance:  100,
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

	n := 1
	transferAmt := float64(10)
	txErrs := make(chan error)
	insertErrs := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3*time.Second))
			//defer cancel()
			ctx := context.WithValue(context.Background(), txKey, txName)
			txErr, insertErr := transferService.TransferFunds(ctx, &models.Transfer{
				SenderID:   senderAccount.ID,
				ReceiverID: receiverAccount.ID,
				Amount:     transferAmt,
			})
			txErrs <- txErr
			insertErrs <- insertErr
		}()
	}

	for i := 0; i < n; i++ {
		txErr := <-txErrs
		insertErr := <-insertErrs
		require.NoError(t, txErr)
		require.NoError(t, insertErr)
	}

	expectedSenderBalance := accounts[0].Balance - float64(n)*transferAmt
	senderAccount, err = accountService.GetAccountByName(testUserID, account1Name)
	require.NoError(t, err)
	require.Equal(t, expectedSenderBalance, senderAccount.Balance, "Sender remaining balance not as expected")
	require.Greater(t, expectedSenderBalance, 0.0)

	expectedRecieverBalance := accounts[1].Balance + float64(n)*transferAmt
	receiverAccount, err = accountService.GetAccountByName(adminUserID, account2Name)
	require.NoError(t, err)
	require.Equal(t, expectedRecieverBalance, receiverAccount.Balance, "Receiver remaining balance not as expected")
	require.Greater(t, expectedRecieverBalance, 100.0)
}
