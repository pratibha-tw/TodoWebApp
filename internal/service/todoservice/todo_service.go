package todoservice

import (
	"todoapp/internal/database/model/todo"
	todorepo "todoapp/internal/repository/todo_repo"
)

type TodoService interface {
	AddTask(t todo.Task) error
	UpdateTask(t todo.Task, userId int) error
	GetTaskById(id int, userId int) (todo.Task, error)
	GetTodoList(id int, criteria todo.TodoSearchCriteria) (todo.Todos, error)
	DeleteTask(id int, userId int) error
	GetTasksNearDueDateButNotCompleted() []todo.Task
	GetNotifications(userId int) []todo.Notification
}
type todoService struct {
	todoRepo todorepo.TodoRepository
}

// GetNotifications implements TodoService.
func (todoService todoService) GetNotifications(userId int) []todo.Notification {
	return todoService.todoRepo.GetNotifications(userId)
}

// GetTasksNearDueDateButNotCompleted implements TodoService.
func (todoService todoService) GetTasksNearDueDateButNotCompleted() []todo.Task {
	return todoService.todoRepo.GetTasksNearDueDateButNotCompleted()
}

// DeleteTask implements TodoService.
func (todoService todoService) DeleteTask(id int, userId int) error {
	return todoService.todoRepo.DeleteTask(id, userId)
}

// GetTodoList implements TodoService.
func (todoService todoService) GetTodoList(id int, criteria todo.TodoSearchCriteria) (todo.Todos, error) {
	return todoService.todoRepo.GetTodoListByUserId(id, criteria)
}

// GetTaskById implements TodoService.
func (todoService todoService) GetTaskById(id int, userId int) (todo.Task, error) {
	return todoService.todoRepo.GetTaskById(id, userId)
}

// UpdateTask implements TodoService.
func (todoService todoService) UpdateTask(t todo.Task, userId int) error {
	return todoService.todoRepo.UpdateTask(t, userId)
}

func (todoService todoService) AddTask(t todo.Task) error {
	return todoService.todoRepo.CreateTask(t)
}

func NewTodoService(todoRepo todorepo.TodoRepository) TodoService {
	return &todoService{todoRepo}
}
