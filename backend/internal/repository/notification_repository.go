package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByID(id uuid.UUID) (*models.Notification, error)
	ListByUser(userID uuid.UUID, page, limit int, unreadOnly bool) ([]models.Notification, int64, error)
	MarkAsRead(id uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	Delete(id uuid.UUID) error
	DeleteAllByUser(userID uuid.UUID) error
	GetUnreadCount(userID uuid.UUID) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByID(id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.Preload("Actor").First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) ListByUser(userID uuid.UUID, page, limit int, unreadOnly bool) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	
	if unreadOnly {
		query = query.Where("read = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Actor").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *notificationRepository) MarkAsRead(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Notification{}, "id = ?", id).Error
}

func (r *notificationRepository) DeleteAllByUser(userID uuid.UUID) error {
	return r.db.Delete(&models.Notification{}, "user_id = ?", userID).Error
}

func (r *notificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error
	return count, err
}
