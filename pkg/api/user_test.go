package api

import (
	"bytes"
	"encoding/json"
	"github.com/JeremyCurmi/simpleBank/pkg/models"
	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	registerEndpoint = "/user/register"
	loginEndpoint    = "/user/login"
)

func TestRegisterUser(t *testing.T) {
	m := setup()
	r := SetupRouter(m)
	t.Run("success", func(t *testing.T) {
		userName := utils.RandomUserName()

		registerBody := models.User{
			UserName: userName,
			Password: "testpassword",
		}
		mockResponse := `{"message":"user registered successfully"}`

		payload, err := json.Marshal(registerBody)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, registerEndpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("duplicate user registration", func(t *testing.T) {
		userName := utils.RandomUserName()

		registerBody := models.User{
			UserName: userName,
			Password: "testpassword",
		}
		mockResponse := `{"error":"duplicate username, use another username"}`

		payload, err := json.Marshal(registerBody)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, registerEndpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// second time
		req, _ = http.NewRequest(http.MethodPost, registerEndpoint, bytes.NewBuffer(payload))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginUser(t *testing.T) {
	m := setup()
	r := SetupRouter(m)
	t.Run("user exists success", func(t *testing.T) {
		userName := utils.RandomUserName()

		registerBody := models.User{
			UserName: userName,
			Password: "testpassword",
		}

		payload, err := json.Marshal(registerBody)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, registerEndpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		req, _ = http.NewRequest(http.MethodPost, loginEndpoint, bytes.NewBuffer(payload))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var actualResponse map[string]string

		responseData, _ := ioutil.ReadAll(w.Body)
		err = json.Unmarshal(responseData, &actualResponse)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, utils.UserLoggedInMessage, actualResponse["message"])
	})

	t.Run("User does not exist", func(t *testing.T) {
		registerBody := models.User{
			UserName: "does-not-exist",
			Password: "testpassword",
		}

		payload, err := json.Marshal(registerBody)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, loginEndpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		mockResponse := `{"error":"user not found"}`

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, http.StatusForbidden, w.Code)
		require.Equal(t, mockResponse, string(responseData))
	})
}
