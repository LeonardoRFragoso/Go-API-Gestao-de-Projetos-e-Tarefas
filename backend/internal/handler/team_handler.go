package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/service"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required,min=1"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type AddTeamMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type AddTeamProjectRequest struct {
	ProjectID string `json:"project_id" binding:"required"`
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	team, err := h.teamService.Create(req.Name, req.Description, req.Color, uid)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, team.ToResponse())
}

func (h *TeamHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	team, members, err := h.teamService.GetByIDWithDetails(id)
	if err != nil {
		NotFoundResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, team.ToDetailResponse(members))
}

func (h *TeamHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	if !h.teamService.CanManageTeam(id, uid) {
		ForbiddenResponse(c, "You don't have permission to update this team")
		return
	}

	team, err := h.teamService.Update(id, req.Name, req.Description, req.Color)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, team.ToResponse())
}

func (h *TeamHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	if err := h.teamService.Delete(id, uid); err != nil {
		if err == service.ErrNotTeamOwner {
			ForbiddenResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

func (h *TeamHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	teams, total, err := h.teamService.ListByUser(uid, page, limit)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	response := make([]models.TeamResponse, 0, len(teams))
	for _, team := range teams {
		response = append(response, team.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *TeamHandler) AddMember(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	var req AddTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	memberID, err := uuid.Parse(req.UserID)
	if err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	role := models.TeamRoleMember
	if req.Role != "" {
		role = models.TeamRole(req.Role)
	}

	userID, _ := c.Get("userID")
	actorID := userID.(uuid.UUID)

	if err := h.teamService.AddMember(teamID, memberID, actorID, role); err != nil {
		if err == service.ErrAlreadyTeamMember {
			ConflictResponse(c, err.Error())
			return
		}
		if err == service.ErrNotTeamOwner {
			ForbiddenResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	memberIDStr := c.Param("memberId")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid member ID")
		return
	}

	userID, _ := c.Get("userID")
	actorID := userID.(uuid.UUID)

	if err := h.teamService.RemoveMember(teamID, memberID, actorID); err != nil {
		if err == service.ErrCannotRemoveOwner {
			ForbiddenResponse(c, err.Error())
			return
		}
		if err == service.ErrNotTeamOwner {
			ForbiddenResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

func (h *TeamHandler) UpdateMemberRole(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	memberIDStr := c.Param("memberId")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid member ID")
		return
	}

	var req UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	userID, _ := c.Get("userID")
	actorID := userID.(uuid.UUID)

	if err := h.teamService.UpdateMemberRole(teamID, memberID, actorID, models.TeamRole(req.Role)); err != nil {
		if err == service.ErrNotTeamOwner {
			ForbiddenResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member role updated successfully"})
}

func (h *TeamHandler) GetMembers(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	members, err := h.teamService.GetMembers(teamID)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	response := make([]models.TeamMemberResponse, 0, len(members))
	for _, m := range members {
		memberResp := models.TeamMemberResponse{
			UserID:   m.UserID,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		}
		if m.User.ID != uuid.Nil {
			userResp := m.User.ToResponse()
			memberResp.User = &userResp
		}
		response = append(response, memberResp)
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *TeamHandler) AddProject(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	var req AddTeamProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	userID, _ := c.Get("userID")
	actorID := userID.(uuid.UUID)

	if err := h.teamService.AddProject(teamID, projectID, actorID); err != nil {
		if err == service.ErrNotTeamOwner {
			ForbiddenResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project added to team successfully"})
}

func (h *TeamHandler) RemoveProject(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	userID, _ := c.Get("userID")
	actorID := userID.(uuid.UUID)

	if err := h.teamService.RemoveProject(teamID, projectID, actorID); err != nil {
		if err == service.ErrNotTeamOwner {
			ForbiddenResponse(c, err.Error())
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project removed from team successfully"})
}

func (h *TeamHandler) GetProjects(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		BadRequestResponse(c, "Invalid team ID")
		return
	}

	projects, err := h.teamService.GetProjects(teamID)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	response := make([]models.ProjectResponse, 0, len(projects))
	for _, p := range projects {
		response = append(response, p.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
