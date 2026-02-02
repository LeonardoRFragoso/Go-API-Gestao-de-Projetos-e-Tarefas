package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/service"
)

type ListHandler struct {
	listService    service.ListService
	boardService   service.BoardService
	projectService service.ProjectService
}

func NewListHandler(listService service.ListService, boardService service.BoardService, projectService service.ProjectService) *ListHandler {
	return &ListHandler{
		listService:    listService,
		boardService:   boardService,
		projectService: projectService,
	}
}

type CreateListRequest struct {
	Name    string `json:"name" binding:"required,min=1"`
	BoardID string `json:"board_id" binding:"required"`
}

type UpdateListRequest struct {
	Name string `json:"name"`
}

type ReorderListsRequest struct {
	Positions map[string]int `json:"positions" binding:"required"`
}

// @Summary Create list
// @Description Create a new list in a board
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateListRequest true "List data"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /lists [post]
func (h *ListHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	boardID, err := uuid.Parse(req.BoardID)
	if err != nil {
		BadRequestResponse(c, "Invalid board ID")
		return
	}

	board, err := h.boardService.GetByID(boardID)
	if err != nil {
		NotFoundResponse(c, "Board not found")
		return
	}

	canAccess, err := h.projectService.CanAccess(board.ProjectID, userID.(uuid.UUID))
	if err != nil || !canAccess {
		ForbiddenResponse(c, "Access denied")
		return
	}

	list, err := h.listService.Create(req.Name, boardID)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	CreatedResponse(c, list.ToResponse())
}

// @Summary Get list
// @Description Get list by ID with tasks
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "List ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /lists/{id} [get]
func (h *ListHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid list ID")
		return
	}

	list, err := h.listService.GetByIDWithTasks(id)
	if err != nil {
		NotFoundResponse(c, "List not found")
		return
	}

	SuccessResponse(c, list.ToResponse())
}

// @Summary Update list
// @Description Update list name
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "List ID"
// @Param request body UpdateListRequest true "List data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /lists/{id} [put]
func (h *ListHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid list ID")
		return
	}

	var req UpdateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	list, err := h.listService.Update(id, req.Name)
	if err != nil {
		if err == service.ErrListNotFound {
			NotFoundResponse(c, "List not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, list.ToResponse())
}

// @Summary Delete list
// @Description Delete list and all its tasks
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "List ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /lists/{id} [delete]
func (h *ListHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid list ID")
		return
	}

	if err := h.listService.Delete(id); err != nil {
		if err == service.ErrListNotFound {
			NotFoundResponse(c, "List not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "List deleted successfully")
}

// @Summary Reorder lists
// @Description Update positions of lists in a board
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param boardId path string true "Board ID"
// @Param request body ReorderListsRequest true "List positions"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /boards/{boardId}/lists/reorder [put]
func (h *ListHandler) Reorder(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid board ID")
		return
	}

	var req ReorderListsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	positions := make(map[uuid.UUID]int)
	for idStr, pos := range req.Positions {
		id, err := uuid.Parse(idStr)
		if err != nil {
			BadRequestResponse(c, "Invalid list ID in positions")
			return
		}
		positions[id] = pos
	}

	if err := h.listService.ReorderLists(boardID, positions); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Lists reordered successfully")
}
