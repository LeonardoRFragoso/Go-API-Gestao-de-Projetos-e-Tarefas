package repository

import (
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type LabelRepository interface {
	Create(label *models.Label) error
	FindByID(id uuid.UUID) (*models.Label, error)
	Update(label *models.Label) error
	Delete(id uuid.UUID) error
	ListByBoard(boardID uuid.UUID) ([]models.Label, error)
}

type labelRepository struct {
	db *gorm.DB
}

func NewLabelRepository(db *gorm.DB) LabelRepository {
	return &labelRepository{db: db}
}

func (r *labelRepository) Create(label *models.Label) error {
	return r.db.Create(label).Error
}

func (r *labelRepository) FindByID(id uuid.UUID) (*models.Label, error) {
	var label models.Label
	err := r.db.First(&label, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *labelRepository) Update(label *models.Label) error {
	return r.db.Save(label).Error
}

func (r *labelRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Label{}, "id = ?", id).Error
}

func (r *labelRepository) ListByBoard(boardID uuid.UUID) ([]models.Label, error) {
	var labels []models.Label
	err := r.db.Where("board_id = ?", boardID).Find(&labels).Error
	return labels, err
}
