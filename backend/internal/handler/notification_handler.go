package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/service"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	unreadOnly := c.DefaultQuery("unread_only", "false") == "true"

	notifications, total, err := h.notificationService.ListByUser(uid, page, limit, unreadOnly)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	response := make([]interface{}, 0, len(notifications))
	for _, n := range notifications {
		response = append(response, n.ToResponse())
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

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	count, err := h.notificationService.GetUnreadCount(uid)
	if err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		BadRequestResponse(c, "Invalid notification ID")
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	if err := h.notificationService.MarkAsRead(id, uid); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	if err := h.notificationService.MarkAllAsRead(uid); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		BadRequestResponse(c, "Invalid notification ID")
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	if err := h.notificationService.Delete(id, uid); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}

func (h *NotificationHandler) DeleteAll(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	if err := h.notificationService.DeleteAll(uid); err != nil {
		InternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications deleted"})
}
