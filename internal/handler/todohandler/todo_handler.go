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
	GetTodoList(ctx *gin.Context)
	DeleteTask(ctx *gin.Context)
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
func (todoHandler todoHandler) GetTodoList(ctx *gin.Context) {
	id := ctx.Param("id")
	UserId, _ := strconv.Atoi(id)

	title, ok1 := ctx.GetQuery("title")
	desc, ok2 := ctx.GetQuery("description")
	priority, ok3 := ctx.GetQuery("priority")
	category, ok4 := ctx.GetQuery("category")
	var criteria todo.TodoSearchCriteria
	if ok1 || ok2 || ok3 || ok4 {
		criteria = todo.TodoSearchCriteria{Title: title, Description: desc, Priority: priority, Category: category}
	}
	response, err := todoHandler.todoService.GetTodoList(UserId, criteria)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
}
func (todoHandler todoHandler) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	Id, _ := strconv.Atoi(id)
	err := todoHandler.todoService.DeleteTask(Id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "Task deleted successfully")
}
func NewTodoHandler(todoService todoservice.TodoService) TodoHandler {
	return &todoHandler{todoService: todoService}
}
