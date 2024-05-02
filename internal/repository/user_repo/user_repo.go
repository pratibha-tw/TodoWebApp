package user_repo

import (
	"database/sql"
	"fmt"
	"log"
	user_model "todoapp/internal/database/model/user"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(u user_model.User) error
	GetUser(username string) (user_model.User, error)
	GetEmailIds(uids []any) map[int]user_model.User
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
	password_hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	// Prepare the INSERT statement
	stmt, err := userRepository.DB.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	// Execute the INSERT statement
	_, err = stmt.Exec(u.Username, u.Email, password_hash)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return err
	}

	return nil
}

func (userRepository userRepository) GetEmailIds(uids []any) map[int]user_model.User {
	user_list := make(map[int]user_model.User)
	commaCnt := len(uids) - 1
	query := "select id,email,username from users where id in ("
	for i := 0; i < len(uids); i++ {
		query += "?"
		if commaCnt > 0 {
			query += ","
		}
		commaCnt--
	}
	query += ");"

	stmt, err := userRepository.DB.Prepare(query)
	if err != nil {
		fmt.Println("error in getting user email")
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(uids...)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	for rows.Next() {
		var id int
		var email string
		var username string
		err = rows.Scan(&id, &email, &username)
		if err != nil {
			fmt.Println(err)
			return user_list
		}
		user_list[id] = user_model.User{Email: email, UserCredentials: user_model.UserCredentials{Username: username}}
	}
	return user_list
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}
