package user

type User struct {
	UserCredentials
	Email string `json:"email" binding:"required"`
}

type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
