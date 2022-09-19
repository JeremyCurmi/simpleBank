package api

import (
	"fmt"
	"github.com/JeremyCurmi/simpleBank/pkg/services"

	"github.com/JeremyCurmi/simpleBank/pkg/config"
	"github.com/JeremyCurmi/simpleBank/pkg/database"
	"github.com/JeremyCurmi/simpleBank/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Manager struct {
	logger          *zap.Logger
	accountsService *services.AccountsService
	userService     *services.UserService
	transferService *services.TransferService
}

func New(logger *zap.Logger, accountsService *services.AccountsService, userService *services.UserService, transferService *services.TransferService) *Manager {
	return &Manager{
		logger:          logger,
		accountsService: accountsService,
		userService:     userService,
		transferService: transferService,
	}
}
func (m *Manager) userRoutes(r *gin.RouterGroup) {
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
func (m *Manager) transfersRoutes(r *gin.RouterGroup) {}
func (m *Manager) InitRoutes(r *gin.Engine) {
	m.userRoutes(r.Group(AuthRoute))
	m.accountsRoutes(r.Group(AccountsRoute))
	m.transfersRoutes(r.Group(TransfersRoute))
}

func InitAPI(logger *zap.Logger, db *sqlx.DB) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	dbClient, err := database.New(logger, *config.DBURL, *config.DBConnMaxLifetime)
	if err != nil {
		return err
	}
	dbConn := dbClient.Conn()
	userService := services.NewUserService(logger, database.NewDBUser(dbConn))
	accountsService := services.NewAccountsService(logger, database.NewDBAccount(dbConn))
	transferService := services.NewTransferService(logger, database.NewDBTransfer(dbConn), accountsService)
	apiManager := New(logger, accountsService, userService, transferService)
	apiManager.InitRoutes(r)

	logger.Info("setting up API on " + URL() + " ⚡️")
	err = r.Run(URL())
	if err != nil {
		return err
	}
	return nil
}

func URL() string {
	return fmt.Sprintf("%s:%s", *config.APIHost, *config.APIPort)
}
