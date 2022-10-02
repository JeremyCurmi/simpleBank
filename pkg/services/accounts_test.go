package services

import (
	"github.com/JeremyCurmi/simpleBank/pkg/database"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	accountName    = "test"
	accountBalance = utils.RandomMoney()
	userCurrency   = utils.RandomCurrency()
	accountService *AccountsService
)

func setupAccounts() *AccountsService {
	dbAccount := database.NewDBAccount(dbConn)
	accountService := NewAccountsService(logger, dbAccount)
	return accountService
}
func TestCreateAccount(t *testing.T) {
	accountService = setupAccounts()
	testAccount := models.Account{
		UserID:   testUserID,
		Name:     accountName,
		Currency: userCurrency,
		Balance:  accountBalance,
	}
	err := accountService.CreateAccount(testAccount)
	require.NoError(t, err)
}
func TestUpdateAccount(t *testing.T) {
	accountService = setupAccounts()
	accountName := utils.RandomAccountName()
	testAccount := models.Account{
		UserID:   testUserID,
		Name:     accountName,
		Currency: userCurrency,
		Balance:  accountBalance,
	}
	err := accountService.CreateAccount(testAccount)
	require.NoError(t, err)

	actualAccount, err := accountService.GetAccountByName(testUserID, accountName)
	require.NoError(t, err)
	require.NotEmpty(t, actualAccount)
	require.Equal(t, testAccount.Balance, actualAccount.Balance)

	var accountID = actualAccount.ID
	t.Run("update account balance", func(t *testing.T) {
		updatedModel := models.Account{
			Name:     accountName,
			Balance:  utils.RandomMoney(),
			Currency: utils.RandomCurrency(),
		}
		err := accountService.UpdateAccountBalance(&updatedModel, accountID)
		require.NoError(t, err)
		actualAccount, err = accountService.GetAccountByName(testUserID, accountName)
		require.NoError(t, err)
		require.Equal(t, updatedModel.Balance, actualAccount.Balance)
	})

	t.Run("update account name", func(t *testing.T) {
		testAccountName := utils.RandomAccountName()
		updatedModel := models.Account{
			Name:     testAccountName,
			Balance:  utils.RandomMoney(),
			Currency: utils.RandomCurrency(),
		}
		err := accountService.UpdateAccountName(&updatedModel, accountID)
		require.NoError(t, err)

		actualAccount, err = accountService.GetAccountByName(testUserID, testAccountName)
		require.NoError(t, err)
		require.NotEmpty(t, actualAccount)
	})
}
func TestGetUserAccounts(t *testing.T) {
	accountService = setupAccounts()
	testAccounts := []models.Account{
		{
			UserID:   adminUserID,
			Name:     utils.RandomAccountName(),
			Currency: userCurrency,
			Balance:  utils.RandomMoney(),
		},
		{
			UserID:   adminUserID,
			Name:     utils.RandomAccountName(),
			Currency: userCurrency,
			Balance:  utils.RandomMoney(),
		},
		{
			UserID:   adminUserID,
			Name:     utils.RandomAccountName(),
			Currency: userCurrency,
			Balance:  utils.RandomMoney(),
		},
	}

	for _, account := range testAccounts {
		err := accountService.CreateAccount(account)
		require.NoError(t, err)
	}
	actualAccounts, err := accountService.GetUserAccounts(adminUserID)
	require.NoError(t, err)
	require.Equal(t, len(testAccounts), len(actualAccounts))
}
func TestGetAccount(t *testing.T) {
	accountService = setupAccounts()
	account := utils.RandomAccountName()
	testAccount := models.Account{
		UserID:   testUserID,
		Name:     account,
		Currency: userCurrency,
		Balance:  accountBalance,
	}
	err := accountService.CreateAccount(testAccount)
	require.NoError(t, err)

	t.Run("account exists", func(t *testing.T) {
		actualAccount, err := accountService.GetAccountByName(testUserID, account)
		require.NoError(t, err)
		require.NotEmpty(t, actualAccount)
	})

	t.Run("account does not exist", func(t *testing.T) {
		actualAccount, err := accountService.GetAccountByName(testUserID, "name-does-not-exist")
		require.Error(t, err)
		require.Empty(t, actualAccount)
	})
}
func TestDeleteAccount(t *testing.T) {
	accountService = setupAccounts()
	accountName := utils.RandomAccountName()
	testAccount := models.Account{
		UserID:   testUserID,
		Name:     accountName,
		Currency: userCurrency,
		Balance:  accountBalance,
	}
	err := accountService.CreateAccount(testAccount)
	require.NoError(t, err)

	account, err := accountService.GetAccountByName(testUserID, accountName)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	err = accountService.DeleteAccount(account.ID)
	require.NoError(t, err)

	deleteAccount, err := accountService.GetAccountByName(testUserID, accountName)
	require.Error(t, err)
	require.Empty(t, deleteAccount)
}
