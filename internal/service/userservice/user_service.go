package userservice

import (
	"errors"
	"log"
	"todoapp/internal/database/model/user"
	"todoapp/internal/repository/user_repo"

	"github.com/dgrijalva/jwt-go"
)

var SECRET_KEY = []byte("gosecretkey")

type UserService interface {
	Register(u user.User) error
	Login(u user.UserCredentials) (string, error)
}

type userService struct {
	userRepo user_repo.UserRepository
}

// Login implements UserService.
func (userService userService) Login(u user.UserCredentials) (string, error) {
	res_user, err := userService.userRepo.GetUser(u.Username)
	if err != nil {
		return "", err
	}
	if res_user.Username != u.Username || res_user.Password != u.Password {
		return "", errors.New("please provide valid username/password")
	}
	token, err := GenerateJWT()
	if err != nil {
		return "", nil
	}
	return token, err

}

func (userService userService) Register(u user.User) error {
	return userService.userRepo.CreateUser(u)
}

func NewUserService(userRepo user_repo.UserRepository) UserService {
	return &userService{userRepo}
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}
