package todohandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoapp/internal/database/model/todo"
	"todoapp/internal/service/todoservice/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoHandlerTest(t *testing.T) {

	t.Run("ShouldReturnStatusCreatedForAddingValidTask", func(t *testing.T) {
		task := todo.Task{Title: "sample task", UserId: 1}
		todoService := mocks.TodoService{}
		todoService.On("AddTask", task).Return(nil)
		todoHandler := NewTodoHandler(&todoService)

		requestBody, err := json.Marshal(task)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		responseRecorder := postResponse(t, todoHandler.AddTask, "/todoapp/api/task/add", "/todoapp/api/task/add", requestBody)
		assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	})
	t.Run("ShouldReturnBadRequestWhenAddingTaskwithEmptyTitle", func(t *testing.T) {
		task := todo.Task{UserId: 1}
		todoService := mocks.TodoService{}
		todoService.On("AddTask", task).Return(errors.New("title can't be empty"))
		todoHandler := NewTodoHandler(&todoService)

		requestBody, err := json.Marshal(task)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		responseRecorder := postResponse(t, todoHandler.AddTask, "/todoapp/api/task/add", "/todoapp/api/task/add", requestBody)
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	})

}

func getResponse(t *testing.T, handlerFunc gin.HandlerFunc, handlerUrl, url string, body []byte) *httptest.ResponseRecorder {
	engine := gin.Default()
	engine.GET(handlerUrl, handlerFunc)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(body))
	require.NoError(t, err)
	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)
	return responseRecorder
}

func postResponse(t *testing.T, handlerFunc gin.HandlerFunc, handlerUrl, url string, body []byte) *httptest.ResponseRecorder {
	engine := gin.Default()
	engine.POST(handlerUrl, handlerFunc)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	require.NoError(t, err)
	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)
	return responseRecorder
}
