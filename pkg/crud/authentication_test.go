package crud

import (
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testUserID       uint   = 99999999
	testPassword     string = "test"
	testPasswordHash string = "$2a$10$Y0FndaIE3OGvkem/AOPKN.ZOKiXYNzvwL533MG8pPH/sIdLkCZWhi"
)

func TestCreateUser(t *testing.T) {
	AuthService := NewAuthService(logger, dbConn)

	newUser := models.User{
		UserName: "test",
		Password: testPassword,
	}

	err := AuthService.CreateUser(&newUser)
	require.NoError(t, err)

	user, err := AuthService.GetUser("test")
	require.NoError(t, err)
	require.NotEqual(t, "test", user.Password)
	require.True(t, utils.CheckPasswordHash("test", user.Password))
}

func TestGetUser(t *testing.T) {
	AuthService := NewAuthService(logger, dbConn)

	user, err := AuthService.GetUser("admin")
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testPasswordHash, user.Password)
	require.Equal(t, testUserID, user.ID)
}

func TestGetUserByID(t *testing.T) {
	AuthService := NewAuthService(logger, dbConn)
	user, err := AuthService.GetUserByID(testUserID)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testPasswordHash, user.Password)
	require.Equal(t, testUserID, user.ID)
}

func TestValidateUser(t *testing.T) {
	AuthService := NewAuthService(logger, dbConn)
	user := models.UserAPI{UserName: "admin", Password: testPassword}
	token, err := AuthService.ValidateUser(&user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}
