package api

import (
	"fmt"

	"github.com/JeremyCurmi/simpleBank/pkg/config"
	"github.com/JeremyCurmi/simpleBank/pkg/crud"
	"github.com/JeremyCurmi/simpleBank/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Manager struct {
	logger          *zap.Logger
	accountsService *crud.AccountsService
	authService     *crud.AuthService
}

func New(logger *zap.Logger, accountsService *crud.AccountsService, authService *crud.AuthService) *Manager {
	return &Manager{logger, accountsService, authService}
}

func (m *Manager) authRoutes(r *gin.RouterGroup) {
	r.POST(AuthLoginRoute, m.login)
	r.POST(AuthRegisterRoute, m.register)
}

func (m *Manager) accountsRoutes(r *gin.RouterGroup) {
	r.Use(middleware.JwtAuthMiddleware())
	r.GET("", m.getAccounts)
	r.GET(AccountsByUser, m.getAccountsByUser)
	r.GET(AccountByAccountID, m.getAccount)
	r.POST("", m.createAccount)
	r.PUT(AccountByAccountID, m.updateAccount)
	r.DELETE(AccountByAccountID, m.deleteAccount)
}

func (m *Manager) InitRoutes(r *gin.Engine) {
	m.authRoutes(r.Group(AuthRoute))
	m.accountsRoutes(r.Group(AccountsRoute))
}

func InitAPI(logger *zap.Logger, db *sqlx.DB) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	authCrud := crud.NewAuthService(logger, db)
	accountsCrud := crud.NewAccountsService(logger, db)
	apiManager := New(logger, accountsCrud, authCrud)
	apiManager.InitRoutes(r)

	logger.Info("setting up API on " + URL() + " ⚡️")
	err := r.Run(URL())
	if err != nil {
		return err
	}
	return nil
}

func URL() string {
	return fmt.Sprintf("%s:%s", *config.APIHost, *config.APIPort)
}