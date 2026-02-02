package repository

import (
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByID(id uuid.UUID) (*models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uuid.UUID) error
	ListByTask(taskID uuid.UUID, page, limit int) ([]models.Comment, int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByID(id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").First(&comment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Comment{}, "id = ?", id).Error
}

func (r *commentRepository) ListByTask(taskID uuid.UUID, page, limit int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Comment{}).Where("task_id = ?", taskID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("User").
		Where("task_id = ?", taskID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
