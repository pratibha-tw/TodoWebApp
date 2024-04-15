package todo

import (
	"time"
)

type Notification struct {
	Id       int       `json:"id" validate:"_"`
	TaskId   int       `json:"task_id" validate:"required"`
	TaskName string    `json:"details" validate:"required"`
	Duedate  time.Time `json:"due_date" validate:"_"`
}
