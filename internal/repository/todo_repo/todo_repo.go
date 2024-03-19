package todorepo

import (
	"database/sql"
	"errors"
	"fmt"
	"todoapp/internal/database/model/todo"
)

type TodoRepository interface {
	CreateTask(t todo.Task) error
	UpdateTask(t todo.Task) error
}

type todoRepository struct {
	db *sql.DB
}

// UpdateTask implements TodoRepository.
func (todorepo todoRepository) UpdateTask(t todo.Task) error {
	query := "UPDATE tasks SET title=?,description=?,priority=?,category=?,due_date=?,done=? where id=?"

	stmt, err := todorepo.db.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Title, t.Description, t.Priority, t.Category, t.Duedate, t.Done, t.ID)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return errors.New("error in updating task")
	}
	fmt.Println(t)
	return nil
}

func (todorepo todoRepository) CreateTask(t todo.Task) error {

	query := "INSERT INTO tasks (title,user_id"
	values := []interface{}{t.Title, t.UserId}
	if t.Description != "" {
		query += ",description"
		values = append(values, t.Description)
	}
	if t.Priority != "" {
		query += ",priority"
		values = append(values, t.Priority)
	}
	if t.Category != "" {
		query += ",category"
		values = append(values, t.Category)
	}
	if !t.Duedate.IsZero() {
		query += ",due_date"
		values = append(values, t.Duedate)
	}
	query += ") values(?,?"

	if t.Description != "" {
		query += ",?"
	}
	if t.Priority != "" {
		query += ",?"
	}
	if t.Category != "" {
		query += ",?"
	}
	if !t.Duedate.IsZero() {
		query += ",?"
	}
	query += ")"
	// Prepare the INSERT statement
	stmt, err := todorepo.db.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	// Execute the INSERT statement
	_, err = stmt.Exec(values...)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return errors.New("error in creating task")
	}

	return nil
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db}
}
