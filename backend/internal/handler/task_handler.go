package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/service"
)

type TaskHandler struct {
	taskService    service.TaskService
	listService    service.ListService
	boardService   service.BoardService
	projectService service.ProjectService
}

func NewTaskHandler(taskService service.TaskService, listService service.ListService, boardService service.BoardService, projectService service.ProjectService) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		listService:    listService,
		boardService:   boardService,
		projectService: projectService,
	}
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=1"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	ListID      string `json:"list_id" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}

type MoveTaskRequest struct {
	ListID   string `json:"list_id" binding:"required"`
	Position int    `json:"position"`
}

type ReorderTasksRequest struct {
	Positions map[string]int `json:"positions" binding:"required"`
}

type AssigneeRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type LabelRequest struct {
	LabelID string `json:"label_id" binding:"required"`
}

// @Summary Create task
// @Description Create a new task in a list
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTaskRequest true "Task data"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	listID, err := uuid.Parse(req.ListID)
	if err != nil {
		BadRequestResponse(c, "Invalid list ID")
		return
	}

	priority := models.TaskPriority(req.Priority)
	if priority == "" {
		priority = models.PriorityMedium
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", req.DueDate)
			if err != nil {
				BadRequestResponse(c, "Invalid due date format")
				return
			}
		}
		dueDate = &parsed
	}

	task, err := h.taskService.Create(req.Title, req.Description, priority, dueDate, listID, userID.(uuid.UUID))
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	CreatedResponse(c, task.ToResponse())
}

// @Summary Get task
// @Description Get task by ID with details
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	task, err := h.taskService.GetByIDWithDetails(id)
	if err != nil {
		NotFoundResponse(c, "Task not found")
		return
	}

	SuccessResponse(c, task.ToResponse())
}

// @Summary Update task
// @Description Update task details
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body UpdateTaskRequest true "Task data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	priority := models.TaskPriority(req.Priority)

	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", req.DueDate)
			if err != nil {
				BadRequestResponse(c, "Invalid due date format")
				return
			}
		}
		dueDate = &parsed
	}

	task, err := h.taskService.Update(id, req.Title, req.Description, priority, dueDate)
	if err != nil {
		if err == service.ErrTaskNotFound {
			NotFoundResponse(c, "Task not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, task.ToResponse())
}

// @Summary Delete task
// @Description Delete task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	if err := h.taskService.Delete(id); err != nil {
		if err == service.ErrTaskNotFound {
			NotFoundResponse(c, "Task not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Task deleted successfully")
}

// @Summary Move task
// @Description Move task to another list or position
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body MoveTaskRequest true "Move data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /tasks/{id}/move [put]
func (h *TaskHandler) Move(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	var req MoveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	listID, err := uuid.Parse(req.ListID)
	if err != nil {
		BadRequestResponse(c, "Invalid list ID")
		return
	}

	if err := h.taskService.MoveTask(id, listID, req.Position); err != nil {
		if err == service.ErrTaskNotFound {
			NotFoundResponse(c, "Task not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	task, _ := h.taskService.GetByIDWithDetails(id)
	SuccessResponse(c, task.ToResponse())
}

// @Summary Reorder tasks
// @Description Update positions of tasks in a list
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param listId path string true "List ID"
// @Param request body ReorderTasksRequest true "Task positions"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /lists/{listId}/tasks/reorder [put]
func (h *TaskHandler) Reorder(c *gin.Context) {
	listID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid list ID")
		return
	}

	var req ReorderTasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	positions := make(map[uuid.UUID]int)
	for idStr, pos := range req.Positions {
		id, err := uuid.Parse(idStr)
		if err != nil {
			BadRequestResponse(c, "Invalid task ID in positions")
			return
		}
		positions[id] = pos
	}

	if err := h.taskService.ReorderTasks(listID, positions); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Tasks reordered successfully")
}

// @Summary Add assignee
// @Description Add assignee to task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body AssigneeRequest true "Assignee data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /tasks/{id}/assignees [post]
func (h *TaskHandler) AddAssignee(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	var req AssigneeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	if err := h.taskService.AddAssignee(taskID, userID); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	task, _ := h.taskService.GetByIDWithDetails(taskID)
	SuccessResponse(c, task.ToResponse())
}

// @Summary Remove assignee
// @Description Remove assignee from task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param userId path string true "User ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /tasks/{id}/assignees/{userId} [delete]
func (h *TaskHandler) RemoveAssignee(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	if err := h.taskService.RemoveAssignee(taskID, userID); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Assignee removed successfully")
}

// @Summary Add label
// @Description Add label to task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body LabelRequest true "Label data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /tasks/{id}/labels [post]
func (h *TaskHandler) AddLabel(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	var req LabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	labelID, err := uuid.Parse(req.LabelID)
	if err != nil {
		BadRequestResponse(c, "Invalid label ID")
		return
	}

	if err := h.taskService.AddLabel(taskID, labelID); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	task, _ := h.taskService.GetByIDWithDetails(taskID)
	SuccessResponse(c, task.ToResponse())
}

// @Summary Remove label
// @Description Remove label from task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param labelId path string true "Label ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /tasks/{id}/labels/{labelId} [delete]
func (h *TaskHandler) RemoveLabel(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	labelID, err := uuid.Parse(c.Param("labelId"))
	if err != nil {
		BadRequestResponse(c, "Invalid label ID")
		return
	}

	if err := h.taskService.RemoveLabel(taskID, labelID); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Label removed successfully")
}

// @Summary Search tasks
// @Description Search tasks in a project
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedResponse
// @Router /projects/{projectId}/tasks/search [get]
func (h *TaskHandler) Search(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	query := c.Query("q")
	if query == "" {
		BadRequestResponse(c, "Search query is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	tasks, total, err := h.taskService.Search(query, projectID, page, limit)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]interface{}, len(tasks))
	for i, task := range tasks {
		responses[i] = task.ToResponse()
	}

	PaginatedSuccessResponse(c, responses, page, limit, total)
}
