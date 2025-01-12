package api

import (
	"net/http"
	"strings"
	"time"

	db "github.com/emzhofb/todo-crud/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required,status"`
	DueDate     string `json:"due_date" binding:"required"`
}

func (server *Server) createTask(ctx *gin.Context) {
	var req createTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dueDate := req.DueDate
	if !strings.Contains(dueDate, "T") {
		dueDate += "T00:00:00Z"
	}

	parsedDueDate, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateTaskParams{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     parsedDueDate,
	}

	task, err := server.store.CreateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}
