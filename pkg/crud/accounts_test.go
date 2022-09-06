package crud

import (
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	accountName    string = "test"
	accountBalance int    = int(utils.RandomMoney())
	userCurrency   string = utils.RandomCurrency()
)

func TestCreateAccount(t *testing.T) {
	AccountService := NewAccountsService(logger, dbConn)
	testAccount := models.Account{
		UserID:   testUserID,
		Name:     accountName,
		Currency: userCurrency,
		Balance:  accountBalance,
	}
	_, err := AccountService.CreateAccount(testAccount)
	require.NoError(t, err)
}

func TestUpdateAccount(t *testing.T) {
	AccountService := NewAccountsService(logger, dbConn)
	accountName := utils.RandomAccountName()
	testAccount := models.Account{
		UserID:   testUserID,
		Name:     accountName,
		Currency: userCurrency,
		Balance:  accountBalance,
	}
	_, err := AccountService.CreateAccount(testAccount)
	if err != nil {
		panic(err)
	}

	actualAccount, err := AccountService.GetAccount(testUserID, accountName)
	require.NoError(t, err)
	require.NotEmpty(t, actualAccount)

	var accountID uint = actualAccount.ID
	t.Run("update account balance", func(t *testing.T) {
		updatedModel := models.Account{
			Name:     accountName,
			Balance:  int(utils.RandomMoney()),
			Currency: utils.RandomCurrency(),
		}
		result, err := AccountService.UpdateUserAccountBalance(accountID, updatedModel)
		require.NoError(t, err)
		require.NotEqual(t, 0, result)
	})

	t.Run("update account name", func(t *testing.T) {
		testAccountName := utils.RandomAccountName()
		updatedModel := models.Account{
			Name:     testAccountName,
			Balance:  int(utils.RandomMoney()),
			Currency: utils.RandomCurrency(),
		}
		result, err := AccountService.UpdateUserAccountBalance(accountID, updatedModel)
		require.NoError(t, err)
		require.Equal(t, 1, result)
	})
}

func TestGetUserAccounts(t *testing.T) {
	AccountService := NewAccountsService(logger, dbConn)
	testAccounts := []models.Account{
		{
			UserID:   adminUserID,
			Name:     utils.RandomAccountName(),
			Currency: userCurrency,
			Balance:  int(utils.RandomMoney()),
		},
		{
			UserID:   adminUserID,
			Name:     utils.RandomAccountName(),
			Currency: userCurrency,
			Balance:  int(utils.RandomMoney()),
		},
		{
			UserID:   adminUserID,
			Name:     utils.RandomAccountName(),
			Currency: userCurrency,
			Balance:  int(utils.RandomMoney()),
		},
	}

	for _, account := range testAccounts {
		_, err := AccountService.CreateAccount(account)
		require.NoError(t, err)
	}
	actualAccounts, err := AccountService.GetUserAccounts(adminUserID)
	require.NoError(t, err)
	require.Equal(t, len(testAccounts), len(actualAccounts))
}

func TestGetAccount(t *testing.T) {
	AccountService := NewAccountsService(logger, dbConn)
	account := utils.RandomAccountName()
	testAccount := models.Account{
		UserID:   testUserID,
		Name:     account,
		Currency: userCurrency,
		Balance:  accountBalance,
	}
	_, err := AccountService.CreateAccount(testAccount)
	require.NoError(t, err)
	t.Run("account exists", func(t *testing.T) {
		actualAccount, err := AccountService.GetAccount(testUserID, account)
		require.NoError(t, err)
		require.NotEmpty(t, actualAccount)
	})

	t.Run("account does not exist", func(t *testing.T) {
		actualAccount, err := AccountService.GetAccount(testUserID, "account-name-does-not-exist")
		require.Error(t, err)
		require.Empty(t, actualAccount)
	})
}

func TestDeleteAccount(t *testing.T) {
	AccountService := NewAccountsService(logger, dbConn)
	accountName := utils.RandomAccountName()
	testAccount := models.Account{
		UserID:   testUserID,
		Name:     accountName,
		Currency: userCurrency,
		Balance:  accountBalance,
	}
	_, err := AccountService.CreateAccount(testAccount)
	require.NoError(t, err)

	res, err := AccountService.DeleteAccount(testUserID, accountName)
	require.NoError(t, err)
	require.Equal(t, 1, res)
}
