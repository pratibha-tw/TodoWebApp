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
	u := user_model.User{UserCredentials: user_model.UserCredentials{Username: "sample_use", Password: "password"}, Email: "sample@gmail.com"}
	t.Run("ShouldRegisterUserWithValidUsername", func(t *testing.T) {

		userService := mocks.UserService{}
		userHandler := NewUserHandler(&userService)
		engine := gin.Default()
		engine.POST("/todoapp/api/register", userHandler.Register)
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
		requestBody, err := json.Marshal(u)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		userService.On("Register", u).Return(expectedErr)
		request, err := http.NewRequest(http.MethodPost, "/todoapp/api/register", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		responseRecorder := httptest.NewRecorder()
		engine.ServeHTTP(responseRecorder, request)
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	})
}
