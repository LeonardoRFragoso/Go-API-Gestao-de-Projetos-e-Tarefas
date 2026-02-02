package repository

import (
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type ListRepository interface {
	Create(list *models.List) error
	FindByID(id uuid.UUID) (*models.List, error)
	FindByIDWithTasks(id uuid.UUID) (*models.List, error)
	Update(list *models.List) error
	Delete(id uuid.UUID) error
	ListByBoard(boardID uuid.UUID) ([]models.List, error)
	GetMaxPosition(boardID uuid.UUID) (int, error)
	UpdatePositions(boardID uuid.UUID, positions map[uuid.UUID]int) error
}

type listRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) ListRepository {
	return &listRepository{db: db}
}

func (r *listRepository) Create(list *models.List) error {
	return r.db.Create(list).Error
}

func (r *listRepository) FindByID(id uuid.UUID) (*models.List, error) {
	var list models.List
	err := r.db.First(&list, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *listRepository) FindByIDWithTasks(id uuid.UUID) (*models.List, error) {
	var list models.List
	err := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Preload("Tasks.Assignees").Preload("Tasks.Labels").First(&list, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *listRepository) Update(list *models.List) error {
	return r.db.Save(list).Error
}

func (r *listRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.List{}, "id = ?", id).Error
}

func (r *listRepository) ListByBoard(boardID uuid.UUID) ([]models.List, error) {
	var lists []models.List
	err := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Where("board_id = ?", boardID).Order("position ASC").Find(&lists).Error
	return lists, err
}

func (r *listRepository) GetMaxPosition(boardID uuid.UUID) (int, error) {
	var maxPos int
	err := r.db.Model(&models.List{}).
		Where("board_id = ?", boardID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPos).Error
	return maxPos, err
}

func (r *listRepository) UpdatePositions(boardID uuid.UUID, positions map[uuid.UUID]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for listID, pos := range positions {
			if err := tx.Model(&models.List{}).
				Where("id = ? AND board_id = ?", listID, boardID).
				Update("position", pos).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
