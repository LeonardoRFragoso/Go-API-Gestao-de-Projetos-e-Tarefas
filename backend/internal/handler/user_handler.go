package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type UpdateUserRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		NotFoundResponse(c, "User not found")
		return
	}

	SuccessResponse(c, user.ToResponse())
}

// @Summary Update user
// @Description Update user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateUserRequest true "User data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /users/me [put]
func (h *UserHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	user, err := h.userService.Update(userID.(uuid.UUID), req.Name, req.Avatar)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, user.ToResponse())
}

// @Summary Update password
// @Description Update user password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdatePasswordRequest true "Password data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /users/me/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	if err := h.userService.UpdatePassword(userID.(uuid.UUID), req.OldPassword, req.NewPassword); err != nil {
		if err == service.ErrInvalidCredentials {
			BadRequestResponse(c, "Invalid old password")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Password updated successfully")
}

// @Summary List users
// @Description List all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedResponse
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userService.List(page, limit)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]interface{}, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	PaginatedSuccessResponse(c, responses, page, limit, total)
}
