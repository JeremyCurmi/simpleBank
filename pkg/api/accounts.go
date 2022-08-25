package api

import (
	"github.com/gin-gonic/gin"
)

const (
	AccountsRoute = "/accounts"
	AccountsByUser = "/user_id/:user_id"
	AccountByAccountID = "account/:account_id"
)

func (m *Manager) getAccounts(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "getAccounts",
	})
}

func (m *Manager) getAccountsByUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "getAccountsByUser",
	})
}

func (m *Manager) getAccount (c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "getAccount",
	})
}

func (m *Manager) createAccount(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "create",
	})
}

func (m *Manager) updateAccount(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "update",
	})
}

func (m *Manager) deleteAccount(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "delete",
	})
}
