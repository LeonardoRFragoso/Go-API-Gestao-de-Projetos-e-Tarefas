package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/service"
)

type BoardHandler struct {
	boardService   service.BoardService
	projectService service.ProjectService
}

func NewBoardHandler(boardService service.BoardService, projectService service.ProjectService) *BoardHandler {
	return &BoardHandler{
		boardService:   boardService,
		projectService: projectService,
	}
}

type CreateBoardRequest struct {
	Name        string `json:"name" binding:"required,min=1"`
	Description string `json:"description"`
	ProjectID   string `json:"project_id" binding:"required"`
}

type UpdateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// @Summary Create board
// @Description Create a new board in a project
// @Tags boards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateBoardRequest true "Board data"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 403 {object} Response
// @Router /boards [post]
func (h *BoardHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		BadRequestResponse(c, "Invalid project ID")
		return
	}

	canAccess, err := h.projectService.CanAccess(projectID, userID.(uuid.UUID))
	if err != nil || !canAccess {
		ForbiddenResponse(c, "Access denied")
		return
	}

	board, err := h.boardService.Create(req.Name, req.Description, projectID)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	CreatedResponse(c, board.ToResponse())
}

// @Summary Get board
// @Description Get board by ID with lists and tasks
// @Tags boards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Board ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /boards/{id} [get]
func (h *BoardHandler) GetByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid board ID")
		return
	}

	board, err := h.boardService.GetByIDWithLists(id)
	if err != nil {
		NotFoundResponse(c, "Board not found")
		return
	}

	canAccess, err := h.projectService.CanAccess(board.ProjectID, userID.(uuid.UUID))
	if err != nil || !canAccess {
		ForbiddenResponse(c, "Access denied")
		return
	}

	listsResponse := make([]interface{}, len(board.Lists))
	for i, list := range board.Lists {
		listsResponse[i] = list.ToResponse()
	}

	SuccessResponse(c, gin.H{
		"id":          board.ID,
		"name":        board.Name,
		"description": board.Description,
		"project_id":  board.ProjectID,
		"lists":       listsResponse,
		"created_at":  board.CreatedAt,
		"updated_at":  board.UpdatedAt,
	})
}

// @Summary Update board
// @Description Update board details
// @Tags boards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Board ID"
// @Param request body UpdateBoardRequest true "Board data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /boards/{id} [put]
func (h *BoardHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid board ID")
		return
	}

	board, err := h.boardService.GetByID(id)
	if err != nil {
		NotFoundResponse(c, "Board not found")
		return
	}

	canAccess, err := h.projectService.CanAccess(board.ProjectID, userID.(uuid.UUID))
	if err != nil || !canAccess {
		ForbiddenResponse(c, "Access denied")
		return
	}

	var req UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	updatedBoard, err := h.boardService.Update(id, req.Name, req.Description)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, updatedBoard.ToResponse())
}

// @Summary Delete board
// @Description Delete board
// @Tags boards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Board ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /boards/{id} [delete]
func (h *BoardHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid board ID")
		return
	}

	board, err := h.boardService.GetByID(id)
	if err != nil {
		NotFoundResponse(c, "Board not found")
		return
	}

	canAccess, err := h.projectService.CanAccess(board.ProjectID, userID.(uuid.UUID))
	if err != nil || !canAccess {
		ForbiddenResponse(c, "Access denied")
		return
	}

	if err := h.boardService.Delete(id); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Board deleted successfully")
}

// @Summary List boards by project
// @Description List all boards in a project
// @Tags boards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Success 200 {object} Response
// @Router /projects/{projectId}/boards [get]
func (h *BoardHandler) ListByProject(c *gin.Context) {
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

	boards, err := h.boardService.ListByProject(projectID)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]interface{}, len(boards))
	for i, board := range boards {
		responses[i] = board.ToResponse()
	}

	SuccessResponse(c, responses)
}
