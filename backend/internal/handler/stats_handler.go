package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/repository"
	"github.com/taskmanager/backend/internal/service"
)

type StatsHandler struct {
	taskRepo       repository.TaskRepository
	projectService service.ProjectService
}

func NewStatsHandler(taskRepo repository.TaskRepository, projectService service.ProjectService) *StatsHandler {
	return &StatsHandler{
		taskRepo:       taskRepo,
		projectService: projectService,
	}
}

type DashboardStatsResponse struct {
	TaskStats     TaskStatsData `json:"task_stats"`
	WeeklyData    []int         `json:"weekly_data"`
	ProjectsCount int64         `json:"projects_count"`
}

type TaskStatsData struct {
	Todo       int `json:"todo"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	uid := userID.(uuid.UUID)

	// Get task statistics
	taskStats, err := h.taskRepo.GetTaskStatsByUser(uid)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to get task statistics")
		return
	}

	// Get weekly completed tasks
	weeklyData, err := h.taskRepo.GetWeeklyCompletedTasks(uid)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to get weekly data")
		return
	}

	// Get projects count
	projects, _, err := h.projectService.ListByUser(uid, 1, 1)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to get projects count")
		return
	}

	projectsCount := int64(0)
	if len(projects) > 0 {
		// Get actual count
		allProjects, _, _ := h.projectService.ListByUser(uid, 1, 1000)
		projectsCount = int64(len(allProjects))
	}

	response := DashboardStatsResponse{
		TaskStats: TaskStatsData{
			Todo:       taskStats["todo"],
			InProgress: taskStats["in_progress"],
			Done:       taskStats["done"],
		},
		WeeklyData:    weeklyData,
		ProjectsCount: projectsCount,
	}

	c.JSON(http.StatusOK, response)
}

func (h *StatsHandler) GetTaskStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	uid := userID.(uuid.UUID)

	taskStats, err := h.taskRepo.GetTaskStatsByUser(uid)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to get task statistics")
		return
	}

	c.JSON(http.StatusOK, TaskStatsData{
		Todo:       taskStats["todo"],
		InProgress: taskStats["in_progress"],
		Done:       taskStats["done"],
	})
}

func (h *StatsHandler) GetWeeklyProgress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	uid := userID.(uuid.UUID)

	weeklyData, err := h.taskRepo.GetWeeklyCompletedTasks(uid)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to get weekly data")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       weeklyData,
		"start_date": time.Now().AddDate(0, 0, -6).Format("2006-01-02"),
		"end_date":   time.Now().Format("2006-01-02"),
	})
}
