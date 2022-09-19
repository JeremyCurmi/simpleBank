package api

import (
	"net/http"

	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	AuthRoute         = "/auth"
	AuthLoginRoute    = "/login"
	AuthRegisterRoute = "/register"
)

func (m *Manager) login(c *gin.Context) {
	var user *models.UserAPI

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := m.userService.ValidateUser(user)
	if err != nil {
		m.logger.Error("Could not validate user: %v", zap.Error(err))
		utils.SendErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	utils.SendOKResponse(c, http.StatusOK, token)
}

func (m *Manager) register(c *gin.Context) {
	var user *models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		m.logger.Error("Could not hash password: %v", zap.Error(err))
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	user.Password = string(hashPassword)

	if err := m.userService.CreateUser(user); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendOKResponse(c, http.StatusCreated, utils.UserCreatedMessage)
}
