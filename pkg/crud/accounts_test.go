package crud

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/go_projects_jer/simple_bank/pkg/models"
	"gitlab.com/go_projects_jer/simple_bank/pkg/utils"
	"go.uber.org/zap"
)

var logger *zap.Logger = utils.NewLogger()

func TestCreateAccount(t *testing.T) {
	AccountService := NewAccountsService(logger, dbConn)

	arg := models.Account{
		UserID:   1,
		Name:     utils.RandomAccountName(),
		Balance:  int(utils.RandomMoney()),
		Currency: utils.RandomCurrency(),
	}

	result, err := AccountService.CreateAccount(arg)
	require.NoError(t, err)
	require.Equal(t, result, 1)
}

func TestUpdateAccount(t *testing.T) {
	AccountService := NewAccountsService(logger, dbConn)

	accountName := utils.RandomAccountName()

	arg := models.Account{
		UserID:   1,
		Name:     accountName,
		Balance:  int(utils.RandomMoney()),
		Currency: utils.RandomCurrency(),
	}
	_, err := AccountService.CreateAccount(arg)
	require.NoError(t, err)

	updatedModel := models.Account{
		Name:     accountName,
		Balance:  int(utils.RandomMoney()),
		Currency: utils.RandomCurrency(),
	}
	_, err = AccountService.UpdateUserAccount(1, accountName, updatedModel)
	require.NoError(t, err)
}
