package todohandler

import (
	"net/http"
	"strconv"
	"todoapp/internal/database/model/todo"
	"todoapp/internal/service/todoservice"

	"github.com/gin-gonic/gin"
)

type TodoHandler interface {
	AddTask(ctx *gin.Context)
	UpdateTask(ctx *gin.Context)
	GetTaskDetails(ctx *gin.Context)
}

type todoHandler struct {
	todoService todoservice.TodoService
}

func (todoHandler todoHandler) AddTask(ctx *gin.Context) {
	var t todo.Task
	if err := ctx.ShouldBindJSON(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := todoHandler.todoService.AddTask(t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusCreated, "Task added successfully")
	}

}
func (todoHandler todoHandler) UpdateTask(ctx *gin.Context) {
	var t todo.Task
	if err := ctx.ShouldBindJSON(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := todoHandler.todoService.UpdateTask(t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusCreated, "Task updated successfully")
	}
}

func (todoHandler todoHandler) GetTaskDetails(ctx *gin.Context) {
	id := ctx.Param("id")
	Id, _ := strconv.Atoi(id)
	response, err := todoHandler.todoService.GetTaskById(Id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func NewTodoHandler(todoService todoservice.TodoService) TodoHandler {
	return &todoHandler{todoService: todoService}
}
