package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewAuthHandler(authService service.AuthService, userService service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 409 {object} Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		if err == service.ErrUserExists {
			ErrorResponse(c, http.StatusConflict, "User already exists")
			return
		}
		InternalServerErrorResponse(c, err.Error())
		return
	}

	CreatedResponse(c, user.ToResponse())
}

// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	tokens, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		UnauthorizedResponse(c, "Invalid credentials")
		return
	}

	user, _ := h.userService.GetByEmail(req.Email)

	SuccessResponse(c, gin.H{
		"user":   user.ToResponse(),
		"tokens": tokens,
	})
}

// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, err.Error())
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		UnauthorizedResponse(c, "Invalid or expired refresh token")
		return
	}

	SuccessResponse(c, tokens)
}

// @Summary Logout user
// @Description Invalidate all refresh tokens for user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if err := h.authService.Logout(userID.(uuid.UUID)); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	MessageResponse(c, "Logged out successfully")
}

// @Summary Get current user
// @Description Get authenticated user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	user, err := h.userService.GetByID(userID.(uuid.UUID))
	if err != nil {
		NotFoundResponse(c, "User not found")
		return
	}

	SuccessResponse(c, user.ToResponse())
}
