package user_repo

import (
	"database/sql"
	"fmt"
	user_model "todoapp/internal/database/model/user"
)

type UserRepository interface {
	CreateUser(u user_model.User) error
	GetUser(username string) (user_model.User, error)
}

type userRepository struct {
	*sql.DB
}

// GetUser implements UserRepository.
func (user_repos userRepository) GetUser(username string) (user_model.User, error) {
	var u user_model.User
	query := fmt.Sprintf("select username,email,password ,id from users where username ='%s'", username)
	row := user_repos.DB.QueryRow(query)
	err := row.Scan(&u.Username, &u.Email, &u.Password, &u.UserId)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (userRepository userRepository) CreateUser(u user_model.User) error {
	// Prepare the INSERT statement
	stmt, err := userRepository.DB.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	// Execute the INSERT statement
	_, err = stmt.Exec(u.Username, u.Email, u.Password)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return err
	}

	return nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}
