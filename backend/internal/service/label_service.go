package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
)

var (
	ErrLabelNotFound = errors.New("label not found")
)

type LabelService interface {
	Create(name, color string, boardID uuid.UUID) (*models.Label, error)
	GetByID(id uuid.UUID) (*models.Label, error)
	Update(id uuid.UUID, name, color string) (*models.Label, error)
	Delete(id uuid.UUID) error
	ListByBoard(boardID uuid.UUID) ([]models.Label, error)
}

type labelService struct {
	labelRepo repository.LabelRepository
}

func NewLabelService(labelRepo repository.LabelRepository) LabelService {
	return &labelService{labelRepo: labelRepo}
}

func (s *labelService) Create(name, color string, boardID uuid.UUID) (*models.Label, error) {
	label := &models.Label{
		Name:    name,
		Color:   color,
		BoardID: boardID,
	}

	if err := s.labelRepo.Create(label); err != nil {
		return nil, err
	}

	return label, nil
}

func (s *labelService) GetByID(id uuid.UUID) (*models.Label, error) {
	label, err := s.labelRepo.FindByID(id)
	if err != nil {
		return nil, ErrLabelNotFound
	}
	return label, nil
}

func (s *labelService) Update(id uuid.UUID, name, color string) (*models.Label, error) {
	label, err := s.labelRepo.FindByID(id)
	if err != nil {
		return nil, ErrLabelNotFound
	}

	if name != "" {
		label.Name = name
	}
	if color != "" {
		label.Color = color
	}

	if err := s.labelRepo.Update(label); err != nil {
		return nil, err
	}

	return label, nil
}

func (s *labelService) Delete(id uuid.UUID) error {
	_, err := s.labelRepo.FindByID(id)
	if err != nil {
		return ErrLabelNotFound
	}
	return s.labelRepo.Delete(id)
}

func (s *labelService) ListByBoard(boardID uuid.UUID) ([]models.Label, error) {
	return s.labelRepo.ListByBoard(boardID)
}
