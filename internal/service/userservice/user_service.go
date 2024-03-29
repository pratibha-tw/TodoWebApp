package userservice

import (
	"errors"
	"log"
	"time"
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
	token, err := GenerateJWT(u)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (userService userService) Register(u user.User) error {
	return userService.userRepo.CreateUser(u)
}

func NewUserService(userRepo user_repo.UserRepository) UserService {
	return &userService{userRepo}
}

func GenerateJWT(u user.UserCredentials) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := user.Claims{Username: u.Username, StandardClaims: jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", errors.New("error in JWT token generation")
	}
	return tokenString, nil
}

func VerifyJWT(tokenString string) error {
	//parse the token
	Claims := user.Claims{}
	_, err := jwt.ParseWithClaims(tokenString, &Claims,
		func(token *jwt.Token) (interface{}, error) {
			return SECRET_KEY, nil
		})

	if err != nil {
		return err
	}
	return nil
}
