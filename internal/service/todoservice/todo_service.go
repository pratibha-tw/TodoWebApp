package todoservice

import (
	"todoapp/internal/database/model/todo"
	todorepo "todoapp/internal/repository/todo_repo"
)

type TodoService interface {
	AddTask(t todo.Task) error
	UpdateTask(t todo.Task) error
}
type todoService struct {
	todoRepo todorepo.TodoRepository
}

// UpdateTask implements TodoService.
func (todoService todoService) UpdateTask(t todo.Task) error {
	return todoService.todoRepo.UpdateTask(t)
}

func (todoService todoService) AddTask(t todo.Task) error {
	return todoService.todoRepo.CreateTask(t)
}

func NewTodoService(todoRepo todorepo.TodoRepository) TodoService {
	return &todoService{todoRepo}
}
