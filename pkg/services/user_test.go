package services

import (
	"github.com/JeremyCurmi/simpleBank/pkg/database"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testPassword     string = "test"
	testPasswordHash string = "$2a$10$Y0FndaIE3OGvkem/AOPKN.ZOKiXYNzvwL533MG8pPH/sIdLkCZWhi"
)

var userService *UserService

func setupUsers() *UserService {
	dbUser := database.NewDBUser(dbConn)
	UserService := NewUserService(logger, dbUser)
	return UserService
}

func TestCreateUser(t *testing.T) {
	userService = setupUsers()
	newUser := models.User{
		UserName: "test",
		Password: testPassword,
	}

	err := userService.CreateUser(&newUser)
	require.NoError(t, err)

	user, err := userService.GetUser("test")
	require.NoError(t, err)
	require.NotEqual(t, "test", user.Password)
	require.True(t, utils.CheckPasswordHash("test", user.Password))
}
func TestGetUser(t *testing.T) {
	userService = setupUsers()

	user, err := userService.GetUser("admin")
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testPasswordHash, user.Password)
	require.Equal(t, testUserID, user.ID)
}
func TestGetUserByID(t *testing.T) {
	userService = setupUsers()
	user, err := userService.GetUserByID(testUserID)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testPasswordHash, user.Password)
	require.Equal(t, testUserID, user.ID)
}
func TestValidateUser(t *testing.T) {
	userService = setupUsers()
	user := models.UserAPI{UserName: "admin", Password: testPassword}
	token, err := userService.ValidateUser(&user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}
