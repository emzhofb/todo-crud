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

type listTodoRequest struct {
	Page  int32 `form:"page"`
	Limit int32 `form:"limit"`
}

func (server *Server) listTask(ctx *gin.Context) {
	var req listTodoRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := int32(5)
	if req.Limit != 0 {
		limit = req.Limit
	}

	page := int32(1)
	if req.Page != 0 {
		page = req.Page
	}
	offset := (page - 1) * limit

	arg := db.ListTasksParams{
		Limit:  limit,
		Offset: offset,
	}

	tasks, err := server.store.ListTasks(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalTasks := len(tasks) // Assuming this function provides the count
	totalPages := (totalTasks + int(limit) - 1) / int(limit)

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_tasks":  totalTasks,
		},
	})
}

type getTaskRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTask(ctx *gin.Context) {
	var req getTaskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	task, err := server.store.GetTask(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}

type updateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required,status"`
	DueDate     string `json:"due_date" binding:"required"`
}

func (server *Server) updateTask(ctx *gin.Context) {
	var query getTaskRequest
	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateTaskRequest
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

	arg := db.UpdateTasksParams{
		ID:          query.ID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     parsedDueDate,
	}

	err = server.store.UpdateTasks(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "task updated successfully",
		"task":    arg,
	})
}

type deleteTaskRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteTask(ctx *gin.Context) {
	var req deleteTaskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteTasks(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}

func (server *Server) createToken(ctx *gin.Context) {
	accessToken, accessPayload, err := server.tokenMaker.CreateToken("admin", "admin", server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires":      accessPayload.ExpiredAt,
	})
}
