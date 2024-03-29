package user

import "github.com/dgrijalva/jwt-go"

type User struct {
	UserCredentials
	Email string `json:"email" binding:"required"`
}

type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
