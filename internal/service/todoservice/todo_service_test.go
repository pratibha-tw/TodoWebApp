package todoservice

import (
	"errors"
	"testing"
	"todoapp/internal/database/model/todo"
	"todoapp/internal/repository/todo_repo/mocks"

	"github.com/stretchr/testify/assert"
)

func TestTodoService(t *testing.T) {

	task := todo.Task{Title: "sample task", UserId: 1}
	updatedTask := todo.Task{ID: 2, Title: "sample task", UserId: 1, Description: "new desc", Done: true}
	t.Run("ShouldBeableToAddTaskForValidUser", func(t *testing.T) {
		mockRepo := mocks.TodoRepository{}
		mockRepo.On("CreateTask", task).Return(nil)
		todoService := NewTodoService(&mockRepo)
		err := todoService.AddTask(task)
		assert.NoError(t, err)
	})

	t.Run("ShouldReturnErrorWhileAddingTaskForInValidUser", func(t *testing.T) {
		mockRepo := mocks.TodoRepository{}
		//actual error msg is different than this msg
		expectedErr := errors.New("user is not present")
		mockRepo.On("CreateTask", task).Return(expectedErr)
		todoService := NewTodoService(&mockRepo)
		err := todoService.AddTask(task)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("ShouldUpdateExistingTask", func(t *testing.T) {
		mockRepo := mocks.TodoRepository{}
		mockRepo.On("UpdateTask", updatedTask).Return(nil)
		todoService := NewTodoService(&mockRepo)
		err := todoService.UpdateTask(updatedTask)
		assert.NoError(t, err)
	})

	t.Run("ShouldReturnErrorForNonExistingTask", func(t *testing.T) {
		mockRepo := mocks.TodoRepository{}
		//actual error msg is different than this msg
		expectedErr := errors.New("task is not present")
		mockRepo.On("UpdateTask", updatedTask).Return(expectedErr)
		todoService := NewTodoService(&mockRepo)
		err := todoService.UpdateTask(updatedTask)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("ShouldReturnTaskDetailsForValidTaskId", func(t *testing.T) {
		mockRepo := mocks.TodoRepository{}
		mockRepo.On("GetTaskById", updatedTask.ID).Return(updatedTask, nil)
		todoService := NewTodoService(&mockRepo)
		res, err := todoService.GetTaskById(2)
		assert.NoError(t, err)
		assert.Equal(t, res.ID, updatedTask.ID)
		assert.Equal(t, res.Title, updatedTask.Title)
	})
	t.Run("ShouldReturnNoDetailsForInValidTaskId", func(t *testing.T) {
		EmptyTask := todo.Task{}
		mockRepo := mocks.TodoRepository{}
		mockRepo.On("GetTaskById", 23).Return(EmptyTask, nil)
		todoService := NewTodoService(&mockRepo)
		res, err := todoService.GetTaskById(23)
		assert.NoError(t, err)
		assert.Equal(t, res, EmptyTask)

	})

	t.Run("ShouldDeleteTaskWithValidTaskId", func(t *testing.T) {
		mockRepo := mocks.TodoRepository{}
		mockRepo.On("DeleteTask", 2).Return(nil)
		todoService := NewTodoService(&mockRepo)
		err := todoService.DeleteTask(2)
		assert.NoError(t, err)
	})
	t.Run("ShouldReturnErrorWhileDeletingTaskWithInValidTaskId", func(t *testing.T) {
		mockRepo := mocks.TodoRepository{}
		expectedErr := errors.New("task does not present")
		mockRepo.On("DeleteTask", 2).Return(expectedErr)
		todoService := NewTodoService(&mockRepo)
		err := todoService.DeleteTask(2)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("ShouldReturnTodoListForGivenUser", func(t *testing.T) {
		mockRes := todo.Todos{TodoList: []todo.Task{task, updatedTask}}
		mockRepo := mocks.TodoRepository{}
		mockRepo.On("GetTodoListByUserId", 1, todo.TodoSearchCriteria{}).Return(mockRes, nil)
		todoService := NewTodoService(&mockRepo)
		res, err := todoService.GetTodoList(1, todo.TodoSearchCriteria{})
		assert.NoError(t, err)
		assert.Equal(t, mockRes, res)
	})
}
