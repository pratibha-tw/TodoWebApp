package todoservice

import (
	"testing"
	"todoapp/internal/database/model/todo"
	"todoapp/internal/repository/todo_repo/mocks"

	"github.com/stretchr/testify/assert"
)

func TestTodoService(t *testing.T) {

	task := todo.Task{Title: "sample task", UserId: 1}
	t.Run("ShouldBeableToAddTaskForValidUser", func(t *testing.T) {
		mockRepo := mocks.TodoRepository{}
		mockRepo.On("CreateTask", task).Return(nil)
		todoService := NewTodoService(&mockRepo)
		err := todoService.AddTask(task)
		assert.NoError(t, err)
	})

}
