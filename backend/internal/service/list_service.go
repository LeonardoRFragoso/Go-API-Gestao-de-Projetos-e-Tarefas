package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
)

var (
	ErrListNotFound = errors.New("list not found")
)

type ListService interface {
	Create(name string, boardID uuid.UUID) (*models.List, error)
	GetByID(id uuid.UUID) (*models.List, error)
	GetByIDWithTasks(id uuid.UUID) (*models.List, error)
	Update(id uuid.UUID, name string) (*models.List, error)
	Delete(id uuid.UUID) error
	ListByBoard(boardID uuid.UUID) ([]models.List, error)
	ReorderLists(boardID uuid.UUID, positions map[uuid.UUID]int) error
}

type listService struct {
	listRepo repository.ListRepository
}

func NewListService(listRepo repository.ListRepository) ListService {
	return &listService{listRepo: listRepo}
}

func (s *listService) Create(name string, boardID uuid.UUID) (*models.List, error) {
	maxPos, err := s.listRepo.GetMaxPosition(boardID)
	if err != nil {
		return nil, err
	}

	list := &models.List{
		Name:     name,
		Position: maxPos + 1,
		BoardID:  boardID,
	}

	if err := s.listRepo.Create(list); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *listService) GetByID(id uuid.UUID) (*models.List, error) {
	list, err := s.listRepo.FindByID(id)
	if err != nil {
		return nil, ErrListNotFound
	}
	return list, nil
}

func (s *listService) GetByIDWithTasks(id uuid.UUID) (*models.List, error) {
	list, err := s.listRepo.FindByIDWithTasks(id)
	if err != nil {
		return nil, ErrListNotFound
	}
	return list, nil
}

func (s *listService) Update(id uuid.UUID, name string) (*models.List, error) {
	list, err := s.listRepo.FindByID(id)
	if err != nil {
		return nil, ErrListNotFound
	}

	if name != "" {
		list.Name = name
	}

	if err := s.listRepo.Update(list); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *listService) Delete(id uuid.UUID) error {
	_, err := s.listRepo.FindByID(id)
	if err != nil {
		return ErrListNotFound
	}
	return s.listRepo.Delete(id)
}

func (s *listService) ListByBoard(boardID uuid.UUID) ([]models.List, error) {
	return s.listRepo.ListByBoard(boardID)
}

func (s *listService) ReorderLists(boardID uuid.UUID, positions map[uuid.UUID]int) error {
	return s.listRepo.UpdatePositions(boardID, positions)
}
