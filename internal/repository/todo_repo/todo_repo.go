package todorepo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"todoapp/internal/database/model/todo"
	errormessages "todoapp/internal/helpers/error_messages"
)

type TodoRepository interface {
	CreateTask(t todo.Task) error
	UpdateTask(t todo.Task, userId int) error
	GetTaskById(id int, userId int) (todo.Task, error)
	GetTodoListByUserId(id int, criteria todo.TodoSearchCriteria) (todo.Todos, error)
	DeleteTask(id int, userId int) error
	GetNotifications(userId int) []todo.Notification
	GetTasksNearDueDateButNotCompleted()
	//AddNotification(wg *sync.WaitGroup, mu *sync.Mutex, data todo.Task) error
	AddNotifications([]todo.Task) error
}

type todoRepository struct {
	db *sql.DB
}

// AddNotifications implements TodoRepository.
func (todorepo todoRepository) AddNotifications(tasks []todo.Task) error {
	count := len(tasks)
	commacnt := count - 1
	values := []any{}
	query := "INSERT INTO notifications (user_id,task_id,details,due_date) values "
	for i := 0; i < count; i++ {
		query += "(?,?,?,?)"
		values = append(values, tasks[i].UserId, tasks[i].ID, tasks[i].Title, tasks[i].Duedate)
		if commacnt > 0 {
			query += ","
		}
		commacnt--
	}
	query += ";"
	_, err := todorepo.db.Exec(query, values...)
	if err != nil {
		fmt.Println("Error executing add notifications query: ", err.Error())
		return err
	}
	return nil
}

func (todorepo todoRepository) GetNotifications(userId int) []todo.Notification {
	var notifications []todo.Notification
	query := "select id,task_id,details,due_date from notifications where user_id=? and created_at >= CURDATE()"
	rows, err := todorepo.db.Query(query, userId)

	if err != nil {
		fmt.Println(err)
		return notifications
	}
	for rows.Next() {
		var n todo.Notification
		err = rows.Scan(&n.Id, &n.TaskId, &n.TaskName, &n.Duedate)
		if err != nil {
			fmt.Println(err)
			return notifications
		}
		notifications = append(notifications, n)
	}
	return notifications
}

// func (todorepo todoRepository) AddNotification(wg *sync.WaitGroup, mu *sync.Mutex, data todo.Task) error {
// 	defer wg.Done()
// 	query := "INSERT INTO notifications (user_id,task_id,details,due_date) values (?,?,?,?)"
// 	mu.Lock()
// 	fmt.Println("Writing to database:", data)
// 	_, err := todorepo.db.Exec(query, data.UserId, data.ID, data.Title, data.Duedate)
// 	if err != nil {
// 		fmt.Println("Error executing add notifications query: ", err.Error())
// 		mu.Unlock()
// 		return err
// 	}
// 	mu.Unlock()
// 	return nil
// }

// GetTasksNearDueDateButNotCompleted implements TodoRepository.
func (todorepo todoRepository) GetTasksNearDueDateButNotCompleted() {
	todos := todo.Todos{}
	query := "select id,title,priority,due_date,user_id from tasks where done=false and DATEDIFF(due_date,NOW()) BETWEEN 0 and 4"
	rows, err := todorepo.db.Query(query)

	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var t todo.Task
		//var due_date sql.NullTime
		err := rows.Scan(&t.ID, &t.Title, &t.Priority, &t.Duedate, &t.UserId)
		if err != nil {
			fmt.Println(err)
		}
		todos.TodoList = append(todos.TodoList, t)
	}
	log.Printf("List of tasks which are due is as follows: %v", todos.TodoList)

	/*var wg sync.WaitGroup
	// Create a Mutex
	var mu sync.Mutex
	for _, task := range todos.TodoList {
		wg.Add(1)
		go todorepo.AddNotification(&wg, &mu, task)
	}*/
	todorepo.AddNotifications(todos.TodoList)
}

// DeleteTask implements TodoRepository.
func (todorepo todoRepository) DeleteTask(id int, userId int) error {
	_, err := todorepo.GetTaskById(id, userId)
	if err != nil {
		switch err.Error() {
		case errormessages.ErrAccessDenied:
			return err
		default:
			return errors.New(errormessages.ErrTaskNotFound)
		}
	}

	query := "DELETE FROM tasks where id=?"
	stmt, err := todorepo.db.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return errors.New("error in deleting task")
	}
	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		return errors.New(errormessages.ErrTaskNotFound)
	}
	return nil
}

// GetTodoListByUserId implements TodoRepository.
func (todorepo todoRepository) GetTodoListByUserId(id int, criteria todo.TodoSearchCriteria) (todo.Todos, error) {
	var list todo.Todos
	query := "select id,title,description,priority,due_date,user_id,category,done from tasks where user_id =?"
	values := []interface{}{id}

	if criteria != (todo.TodoSearchCriteria{}) {
		query += " and (title like ? or description like ? or priority = ? or category = ?)"
		criteria.Title = strings.ReplaceAll(criteria.Title, "\"", "")
		criteria.Description = strings.ReplaceAll(criteria.Description, "\"", "")
		criteria.Priority = strings.ReplaceAll(criteria.Priority, "\"", "")
		criteria.Category = strings.ReplaceAll(criteria.Category, "\"", "")

		if criteria.Title != "" {
			values = append(values, "%"+criteria.Title+"%")
		} else {
			values = append(values, criteria.Title)
		}
		if criteria.Description != "" {
			values = append(values, "%"+criteria.Description+"%")
		} else {
			values = append(values, criteria.Description)
		}
		values = append(values, criteria.Priority, criteria.Category)
	}
	for _, val := range values {
		query = strings.Replace(query, "?", fmt.Sprintf("'%v'", val), 1)
	}
	rows, err := todorepo.db.Query(query)

	if err != nil {
		fmt.Println(err)
		return list, errors.New(errormessages.ErrInvallidUserId)
	}

	for rows.Next() {
		var t todo.Task
		var description sql.NullString
		var due_date sql.NullTime
		var category sql.NullString
		err := rows.Scan(&t.ID, &t.Title, &description, &t.Priority, &due_date, &t.UserId, &category, &t.Done)
		if err != nil {
			fmt.Println(err)
			return list, errors.New("error while getting task details")
		}
		if category.Valid {
			t.Category = category.String
		}
		if description.Valid {
			t.Description = description.String
		}
		if due_date.Valid {
			t.Duedate = due_date.Time
		}
		list.TodoList = append(list.TodoList, t)
	}
	return list, nil
}

// GetTaskById implements TodoRepository.
func (todorepo todoRepository) GetTaskById(id int, userId int) (todo.Task, error) {
	var t todo.Task
	var description sql.NullString
	var due_date sql.NullTime
	var category sql.NullString
	row := todorepo.db.QueryRow("select id,title,description,priority,due_date,user_id,category,done from tasks where id =?", id)
	err := row.Scan(&t.ID, &t.Title, &description, &t.Priority, &due_date, &t.UserId, &category, &t.Done)
	if err != nil {
		fmt.Println(err)
		return t, errors.New(errormessages.ErrInvalidTaskId)
	}
	if userId != t.UserId {
		return t, errors.New(errormessages.ErrAccessDenied)
	}
	if category.Valid {
		t.Category = category.String
	}
	if description.Valid {
		t.Description = description.String
	}
	if due_date.Valid {
		t.Duedate = due_date.Time
	}
	return t, nil
}

// UpdateTask implements TodoRepository.
func (todorepo todoRepository) UpdateTask(t todo.Task, userId int) error {
	_, err := todorepo.GetTaskById(t.ID, userId)
	if err != nil {
		switch err.Error() {
		case errormessages.ErrAccessDenied:
			return err
		default:
			return errors.New(errormessages.ErrTaskNotFound)
		}
	}
	query := "UPDATE tasks SET title=?,description=?,priority=?,category=?,due_date=?,done=? where id=?"

	stmt, err := todorepo.db.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Title, t.Description, t.Priority, t.Category, t.Duedate, t.Done, t.ID)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return errors.New(errormessages.ErrTaskUpdate)
	}
	rowAffected, _ := res.RowsAffected()

	if rowAffected == 0 {
		return errors.New(errormessages.ErrTaskNotFound)
	}
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
		fmt.Println("Error executing query:")
		return err
	}

	return nil
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db}
}
