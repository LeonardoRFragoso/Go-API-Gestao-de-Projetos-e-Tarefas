package repository

import (
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type BoardRepository interface {
	Create(board *models.Board) error
	FindByID(id uuid.UUID) (*models.Board, error)
	FindByIDWithLists(id uuid.UUID) (*models.Board, error)
	Update(board *models.Board) error
	Delete(id uuid.UUID) error
	ListByProject(projectID uuid.UUID) ([]models.Board, error)
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(board *models.Board) error {
	return r.db.Create(board).Error
}

func (r *boardRepository) FindByID(id uuid.UUID) (*models.Board, error) {
	var board models.Board
	err := r.db.Preload("Project").First(&board, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) FindByIDWithLists(id uuid.UUID) (*models.Board, error) {
	var board models.Board
	err := r.db.Preload("Project").
		Preload("Lists", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Lists.Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Lists.Tasks.Assignees").
		Preload("Lists.Tasks.Labels").
		Preload("Lists.Tasks.Creator").
		First(&board, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) Update(board *models.Board) error {
	return r.db.Save(board).Error
}

func (r *boardRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Board{}, "id = ?", id).Error
}

func (r *boardRepository) ListByProject(projectID uuid.UUID) ([]models.Board, error) {
	var boards []models.Board
	err := r.db.Preload("Lists").Where("project_id = ?", projectID).Find(&boards).Error
	return boards, err
}
