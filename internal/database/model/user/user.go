package user

import "github.com/dgrijalva/jwt-go"

type User struct {
	UserCredentials
	Email string `json:"email" binding:"required"`
}

type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserId   int    `json:"user_id"`
}

type Claims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}
