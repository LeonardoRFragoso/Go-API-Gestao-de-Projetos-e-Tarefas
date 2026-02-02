package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/service"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
	TaskID  string `json:"task_id" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

func (h *CommentHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	taskID, err := uuid.Parse(req.TaskID)
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	comment, err := h.commentService.Create(req.Content, taskID, userID.(uuid.UUID))
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	CreatedResponse(c, comment.ToResponse())
}

func (h *CommentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid comment ID")
		return
	}

	comment, err := h.commentService.GetByID(id)
	if err != nil {
		NotFoundResponse(c, "Comment not found")
		return
	}

	SuccessResponse(c, comment.ToResponse())
}

func (h *CommentHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid comment ID")
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	comment, err := h.commentService.Update(id, userID.(uuid.UUID), req.Content)
	if err != nil {
		if err == service.ErrNotCommentOwner {
			ForbiddenResponse(c, "Only the owner can update the comment")
			return
		}
		if err == service.ErrCommentNotFound {
			NotFoundResponse(c, "Comment not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, comment.ToResponse())
}

func (h *CommentHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid comment ID")
		return
	}

	if err := h.commentService.Delete(id, userID.(uuid.UUID)); err != nil {
		if err == service.ErrNotCommentOwner {
			ForbiddenResponse(c, "Only the owner can delete the comment")
			return
		}
		if err == service.ErrCommentNotFound {
			NotFoundResponse(c, "Comment not found")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Comment deleted successfully")
}

func (h *CommentHandler) ListByTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid task ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	comments, total, err := h.commentService.ListByTask(taskID, page, limit)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]interface{}, len(comments))
	for i, comment := range comments {
		responses[i] = comment.ToResponse()
	}

	PaginatedSuccessResponse(c, responses, page, limit, total)
}
