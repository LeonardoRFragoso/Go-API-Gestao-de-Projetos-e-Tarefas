package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/service"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role"`
}

// @Summary Create project
// @Description Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProjectRequest true "Project data"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	if req.Color == "" {
		req.Color = "#3B82F6"
	}

	project, err := h.projectService.Create(req.Name, req.Description, req.Color, userID.(uuid.UUID))
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	CreatedResponse(c, project.ToResponse())
}

// @Summary Get project
// @Description Get project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	canAccess, err := h.projectService.CanAccess(id, userID.(uuid.UUID))
	if err != nil || !canAccess {
		ForbiddenResponse(c, "Access denied")
		return
	}

	project, err := h.projectService.GetByIDWithDetails(id)
	if err != nil {
		NotFoundResponse(c, "Project not found")
		return
	}

	SuccessResponse(c, project.ToResponse())
}

// @Summary Update project
// @Description Update project details
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body UpdateProjectRequest true "Project data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 403 {object} Response
// @Router /projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	project, err := h.projectService.Update(id, userID.(uuid.UUID), req.Name, req.Description, req.Color)
	if err != nil {
		if err == service.ErrNotProjectOwner {
			ForbiddenResponse(c, "Only the owner can update the project")
			return
		}
		if err == service.ErrProjectNotFound {
			NotFoundResponse(c, "Project not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, project.ToResponse())
}

// @Summary Delete project
// @Description Delete project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} Response
// @Failure 403 {object} Response
// @Failure 404 {object} Response
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	if err := h.projectService.Delete(id, userID.(uuid.UUID)); err != nil {
		if err == service.ErrNotProjectOwner {
			ForbiddenResponse(c, "Only the owner can delete the project")
			return
		}
		if err == service.ErrProjectNotFound {
			NotFoundResponse(c, "Project not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Project deleted successfully")
}

// @Summary List user projects
// @Description List projects for authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedResponse
// @Router /projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
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

	projects, total, err := h.projectService.ListByUser(userID.(uuid.UUID), page, limit)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]interface{}, len(projects))
	for i, project := range projects {
		responses[i] = project.ToResponse()
	}

	PaginatedSuccessResponse(c, responses, page, limit, total)
}

// @Summary Add member to project
// @Description Add a member to the project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body AddMemberRequest true "Member data"
// @Success 200 {object} Response
// @Failure 403 {object} Response
// @Router /projects/{id}/members [post]
func (h *ProjectHandler) AddMember(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	memberID, err := uuid.Parse(req.UserID)
	if err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	role := req.Role
	if role == "" {
		role = "member"
	}

	if err := h.projectService.AddMember(projectID, userID.(uuid.UUID), memberID, role); err != nil {
		if err == service.ErrNotProjectOwner {
			ForbiddenResponse(c, "Only the owner can add members")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Member added successfully")
}

// @Summary Remove member from project
// @Description Remove a member from the project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param memberId path string true "Member ID"
// @Success 200 {object} Response
// @Failure 403 {object} Response
// @Router /projects/{id}/members/{memberId} [delete]
func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		BadRequestResponse(c, "Invalid member ID")
		return
	}

	if err := h.projectService.RemoveMember(projectID, userID.(uuid.UUID), memberID); err != nil {
		if err == service.ErrNotProjectOwner {
			ForbiddenResponse(c, "Only the owner can remove members")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Member removed successfully")
}

// @Summary Get project members
// @Description Get all members of a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} Response
// @Router /projects/{id}/members [get]
func (h *ProjectHandler) GetMembers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	canAccess, err := h.projectService.CanAccess(projectID, userID.(uuid.UUID))
	if err != nil || !canAccess {
		ForbiddenResponse(c, "Access denied")
		return
	}

	members, err := h.projectService.GetMembers(projectID)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]interface{}, len(members))
	for i, member := range members {
		responses[i] = member.ToResponse()
	}

	SuccessResponse(c, responses)
}
