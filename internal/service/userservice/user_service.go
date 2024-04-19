package userservice

import (
	"errors"
	"log"
	"time"
	"todoapp/internal/database/model/user"
	"todoapp/internal/repository/user_repo"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
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
	if res_user.Username != u.Username {
		return "", errors.New("please provide valid username")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(res_user.Password), []byte(u.Password)); err != nil {
		return "", errors.New("please provide valid password")
	}
	u.UserId = res_user.UserId
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
	claims := user.Claims{UserId: u.UserId, StandardClaims: jwt.StandardClaims{
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

func VerifyJWT(tokenString string) (user.Claims, error) {
	//parse the token
	Claims := user.Claims{}
	_, err := jwt.ParseWithClaims(tokenString, &Claims,
		func(token *jwt.Token) (interface{}, error) {
			return SECRET_KEY, nil
		})

	if err != nil {
		return Claims, err
	}
	return Claims, nil
}
