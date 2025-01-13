package api

import (
	"fmt"
	"net/http"
	"strconv"
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

type TaskSchema struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

type listTasksRequest struct {
	Status string `form:"status"`
	Search string `form:"search"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

func (server *Server) listTask(ctx *gin.Context) {
	// Parse query parameters
	var req listTasksRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default pagination values
	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}
	page := 1
	if req.Page > 0 {
		page = req.Page
	}
	offset := (page - 1) * limit

	// Base query
	query := "SELECT id, title, description, status, due_date FROM tasks WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM tasks WHERE 1=1"

	// Dynamic filters
	var args []any
	argIndex := 1 // For dynamic parameter binding

	if req.Status != "" {
		query += " AND status = $" + fmt.Sprintf("%d", argIndex)
		countQuery += " AND status = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.Status)
		argIndex++
	}

	if req.Search != "" {
		query += " AND (title ILIKE '%' || $" + fmt.Sprintf("%d", argIndex) + " || '%' OR description ILIKE '%' || $" + fmt.Sprintf("%d", argIndex) + " || '%')"
		countQuery += " AND (title ILIKE '%' || $" + fmt.Sprintf("%d", argIndex) + " || '%' OR description ILIKE '%' || $" + fmt.Sprintf("%d", argIndex) + " || '%')"
		args = append(args, req.Search)
		argIndex++
	}

	// Fetch total task count for pagination
	var totalTasks int
	err := server.connPool.QueryRow(ctx.Request.Context(), countQuery, args...).Scan(&totalTasks)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalPages := (totalTasks + int(limit) - 1) / int(limit)

	// Add pagination
	query += " ORDER BY id LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, strconv.Itoa(limit), strconv.Itoa(offset))

	// Fetch tasks
	rows, err := server.connPool.Query(ctx.Request.Context(), query, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Parse rows into tasks
	var tasks []TaskSchema
	for rows.Next() {
		var task TaskSchema
		var dueDate time.Time
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &dueDate); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		task.DueDate = dueDate.Format("2006-01-02")
		tasks = append(tasks, task)
	}

	// Check for errors in rows iteration
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with tasks and pagination metadata
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

	returnedTask := TaskSchema{
		ID:          int32(task.ID),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate.Format("2006-01-02"),
	}

	ctx.JSON(http.StatusOK, returnedTask)
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

	returnedTask := TaskSchema{
		ID:          int32(query.ID),
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     req.DueDate,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "task updated successfully",
		"task":    returnedTask,
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
