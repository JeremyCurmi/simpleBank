package api

import (
	"github.com/lib/pq"
	"net/http"

	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	UserRoute         = "/user"
	AuthLoginRoute    = "/login"
	AuthRegisterRoute = "/register"
)

func (m *Manager) login(c *gin.Context) {
	var user *models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := m.userService.ValidateUser(user)
	if err != nil {
		m.logger.Error("Could not validate user: %v", zap.Error(err))
		if err.Error() == utils.PqNoRowsFound {
			utils.SendErrorResponse(c, http.StatusForbidden, utils.UserNotFoundMessage)
			return
		}
		utils.SendErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	utils.SendOKResponse(c, http.StatusOK, utils.TokenResponse(utils.UserLoggedInMessage, token))
}
func (m *Manager) register(c *gin.Context) {
	var user *models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := m.userService.CreateUser(user); err != nil {
		switch err.(type) {
		case *pq.Error:
			if err.(*pq.Error).Code == utils.PqDuplicateErrorCode {
				utils.SendErrorResponse(c, http.StatusBadRequest, utils.UserCreateDuplicateMessage)
			}
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendOKResponse(c, http.StatusCreated, utils.MessageResponse(utils.UserCreatedMessage))
}
