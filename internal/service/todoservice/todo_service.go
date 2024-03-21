package todoservice

import (
	"todoapp/internal/database/model/todo"
	todorepo "todoapp/internal/repository/todo_repo"
)

type TodoService interface {
	AddTask(t todo.Task) error
	UpdateTask(t todo.Task) error
	GetTaskById(id int) (todo.Task, error)
	GetTodoList(id int) (todo.Todos, error)
	DeleteTask(id int) error
}
type todoService struct {
	todoRepo todorepo.TodoRepository
}

// DeleteTask implements TodoService.
func (todoService todoService) DeleteTask(id int) error {
	return todoService.todoRepo.DeleteTask(id)
}

// GetTodoList implements TodoService.
func (todoService todoService) GetTodoList(id int) (todo.Todos, error) {
	return todoService.todoRepo.GetTodoListByUserId(id)
}

// GetTaskById implements TodoService.
func (todoService todoService) GetTaskById(id int) (todo.Task, error) {
	return todoService.todoRepo.GetTaskById(id)
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
