package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoapp/internal/service/userservice/mocks"

	user_model "todoapp/internal/database/model/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler(t *testing.T) {
	t.Run("ShouldRegisterUserWithValidUsername", func(t *testing.T) {

		userService := mocks.UserService{}
		userHandler := NewUserHandler(&userService)
		engine := gin.Default()
		engine.POST("/todoapp/api/register", userHandler.Register)
		u := user_model.User{Username: "testuser", Email: "testemail.com", Password: "password"}
		requestBody, err := json.Marshal(u)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		userService.On("Register", u).Return(nil)
		request, err := http.NewRequest(http.MethodPost, "/todoapp/api/register", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		responseRecorder := httptest.NewRecorder()
		engine.ServeHTTP(responseRecorder, request)
		assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	})

	t.Run("ShouldReturnErrorForDuplicateUsernameRegistration", func(t *testing.T) {

		userService := mocks.UserService{}
		userHandler := NewUserHandler(&userService)
		engine := gin.Default()
		engine.POST("/todoapp/api/register", userHandler.Register)
		expectedErr := errors.New("error in creating user")

		u2 := user_model.User{Username: "testuser", Email: "testemail1.com", Password: "password1"}
		requestBody, err := json.Marshal(u2)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		userService.On("Register", u2).Return(expectedErr)
		request, err := http.NewRequest(http.MethodPost, "/todoapp/api/register", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		responseRecorder := httptest.NewRecorder()
		engine.ServeHTTP(responseRecorder, request)
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	})
}
