package todohandler

import (
	"net/http"
	"strconv"
	"todoapp/internal/database/model/todo"
	errormessages "todoapp/internal/helpers/error_messages"
	"todoapp/internal/service/todoservice"

	"github.com/gin-gonic/gin"
)

type TodoHandler interface {
	AddTask(ctx *gin.Context)
	UpdateTask(ctx *gin.Context)
	GetTaskDetails(ctx *gin.Context)
	GetTodoList(ctx *gin.Context)
	DeleteTask(ctx *gin.Context)
	GetNotifications(ctx *gin.Context)
}

type todoHandler struct {
	todoService todoservice.TodoService
}

func (todoHandler todoHandler) AddTask(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")
	var t todo.Task
	if err := ctx.ShouldBindJSON(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userId != t.UserId {
		ctx.JSON(http.StatusForbidden, errormessages.ErrAccessDenied)
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
	userId := ctx.GetInt("user_id")
	var t todo.Task
	if err := ctx.ShouldBindJSON(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := todoHandler.todoService.UpdateTask(t, userId); err != nil {

		switch err.Error() {
		case errormessages.ErrAccessDenied:
			ctx.JSON(http.StatusForbidden, err.Error())
		default:
			ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return

	} else {
		ctx.JSON(http.StatusOK, "Task updated successfully")
	}
}

func (todoHandler todoHandler) GetTaskDetails(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")
	id := ctx.Param("id")
	Id, _ := strconv.Atoi(id)
	response, err := todoHandler.todoService.GetTaskById(Id, userId)

	if err != nil {
		switch err.Error() {
		case errormessages.ErrAccessDenied:
			ctx.JSON(http.StatusForbidden, err.Error())
		default:
			ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)
}
func (todoHandler todoHandler) GetTodoList(ctx *gin.Context) {
	loggedInUserId := ctx.GetInt("user_id")
	id := ctx.Param("id")
	UserId, _ := strconv.Atoi(id)

	if loggedInUserId != UserId {
		ctx.JSON(http.StatusForbidden, errormessages.ErrAccessDenied)
		return
	}
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
	userId := ctx.GetInt("user_id")
	id := ctx.Param("id")
	taskId, _ := strconv.Atoi(id)
	err := todoHandler.todoService.DeleteTask(taskId, userId)

	if err != nil {
		switch err.Error() {
		case errormessages.ErrAccessDenied:
			ctx.JSON(http.StatusForbidden, err.Error())
		default:
			ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, "Task deleted successfully")
}

func (todoHandler todoHandler) GetNotifications(ctx *gin.Context) {
	loggedInUserId := ctx.GetInt("user_id")
	id := ctx.Param("id")
	UserId, _ := strconv.Atoi(id)

	if loggedInUserId != UserId {
		ctx.JSON(http.StatusForbidden, errormessages.ErrAccessDenied)
		return
	}
	response := todoHandler.todoService.GetNotifications(UserId)

	ctx.JSON(http.StatusOK, response)
}
func NewTodoHandler(todoService todoservice.TodoService) TodoHandler {
	return &todoHandler{todoService: todoService}
}
