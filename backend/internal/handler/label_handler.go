package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/service"
)

type LabelHandler struct {
	labelService service.LabelService
}

func NewLabelHandler(labelService service.LabelService) *LabelHandler {
	return &LabelHandler{labelService: labelService}
}

type CreateLabelRequest struct {
	Name    string `json:"name" binding:"required,min=1"`
	Color   string `json:"color"`
	BoardID string `json:"board_id" binding:"required"`
}

type UpdateLabelRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (h *LabelHandler) Create(c *gin.Context) {
	var req CreateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	boardID, err := uuid.Parse(req.BoardID)
	if err != nil {
		BadRequestResponse(c, "Invalid board ID")
		return
	}

	if req.Color == "" {
		req.Color = "#6B7280"
	}

	label, err := h.labelService.Create(req.Name, req.Color, boardID)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	CreatedResponse(c, label.ToResponse())
}

func (h *LabelHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid label ID")
		return
	}

	label, err := h.labelService.GetByID(id)
	if err != nil {
		NotFoundResponse(c, "Label not found")
		return
	}

	SuccessResponse(c, label.ToResponse())
}

func (h *LabelHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid label ID")
		return
	}

	var req UpdateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	label, err := h.labelService.Update(id, req.Name, req.Color)
	if err != nil {
		if err == service.ErrLabelNotFound {
			NotFoundResponse(c, "Label not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, label.ToResponse())
}

func (h *LabelHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid label ID")
		return
	}

	if err := h.labelService.Delete(id); err != nil {
		if err == service.ErrLabelNotFound {
			NotFoundResponse(c, "Label not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Label deleted successfully")
}

func (h *LabelHandler) ListByBoard(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid board ID")
		return
	}

	labels, err := h.labelService.ListByBoard(boardID)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]interface{}, len(labels))
	for i, label := range labels {
		responses[i] = label.ToResponse()
	}

	SuccessResponse(c, responses)
}
