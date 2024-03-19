package todo

import (
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duedate     time.Time `json:"due_date"`
	Priority    string    `json:"priority"`
	Category    string    `json:"category"`
	UserId      int       `json:"user_id"`
	Done        bool      `json:"is_done"`
}

type Todos struct {
	TodoList []Task `json:"todos"`
}
