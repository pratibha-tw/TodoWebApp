package userservice

import (
	"errors"
	"testing"
	user_model "todoapp/internal/database/model/user"
	"todoapp/internal/repository/user_repo/mocks"

	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	t.Run("ShouldRegisterUserWithUniqueUsername", func(t *testing.T) {
		mockUserRepo := mocks.UserRepository{}
		u := user_model.User{Username: "sample_user", Email: "sample@gmail.com", Password: "password"}
		mockUserRepo.On("CreateUser", u).Return(nil)
		userService := NewUserService(&mockUserRepo)
		err := userService.Register(u)
		assert.NoError(t, err)
	})
	t.Run("ShouldNotRegisterUserWithDuplicateUsername", func(t *testing.T) {
		mockUserRepo := mocks.UserRepository{}
		u1 := user_model.User{Username: "sample_user", Email: "sample@gmail.com", Password: "password"}
		u2 := user_model.User{Username: "sample_user", Email: "sample@gmail.com", Password: "password"}
		expectedErr := errors.New("error in creating user")
		mockUserRepo.On("CreateUser", u2).Return(expectedErr)
		userService := NewUserService(&mockUserRepo)
		userService.Register(u1)
		err := userService.Register(u2)
		assert.Equal(t, expectedErr, err)
	})

}
