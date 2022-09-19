package api

import "github.com/gin-gonic/gin"

const (
	TransfersRoute = "/transfers"
)

func (m *Manager) createTransfer(c *gin.Context) {

	// TODO: handle transfer request -> get status success/failed
	// TODO: create transfer record based on 👆status (if error is nil then success)
	// TODO: transfers.Status = if err == nil { success } else { failed }
	c.JSON(200, gin.H{
		"message": "create",
	})
}
