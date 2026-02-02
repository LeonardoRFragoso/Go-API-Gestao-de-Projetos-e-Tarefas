package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
)

var (
	ErrBoardNotFound = errors.New("board not found")
)

type BoardService interface {
	Create(name, description string, projectID uuid.UUID) (*models.Board, error)
	GetByID(id uuid.UUID) (*models.Board, error)
	GetByIDWithLists(id uuid.UUID) (*models.Board, error)
	Update(id uuid.UUID, name, description string) (*models.Board, error)
	Delete(id uuid.UUID) error
	ListByProject(projectID uuid.UUID) ([]models.Board, error)
}

type boardService struct {
	boardRepo   repository.BoardRepository
	listRepo    repository.ListRepository
}

func NewBoardService(boardRepo repository.BoardRepository, listRepo repository.ListRepository) BoardService {
	return &boardService{
		boardRepo: boardRepo,
		listRepo:  listRepo,
	}
}

func (s *boardService) Create(name, description string, projectID uuid.UUID) (*models.Board, error) {
	board := &models.Board{
		Name:        name,
		Description: description,
		ProjectID:   projectID,
	}

	if err := s.boardRepo.Create(board); err != nil {
		return nil, err
	}

	defaultLists := []string{"To Do", "In Progress", "Done"}
	for i, listName := range defaultLists {
		list := &models.List{
			Name:     listName,
			Position: i,
			BoardID:  board.ID,
		}
		if err := s.listRepo.Create(list); err != nil {
			return nil, err
		}
	}

	return s.boardRepo.FindByIDWithLists(board.ID)
}

func (s *boardService) GetByID(id uuid.UUID) (*models.Board, error) {
	board, err := s.boardRepo.FindByID(id)
	if err != nil {
		return nil, ErrBoardNotFound
	}
	return board, nil
}

func (s *boardService) GetByIDWithLists(id uuid.UUID) (*models.Board, error) {
	board, err := s.boardRepo.FindByIDWithLists(id)
	if err != nil {
		return nil, ErrBoardNotFound
	}
	return board, nil
}

func (s *boardService) Update(id uuid.UUID, name, description string) (*models.Board, error) {
	board, err := s.boardRepo.FindByID(id)
	if err != nil {
		return nil, ErrBoardNotFound
	}

	if name != "" {
		board.Name = name
	}
	if description != "" {
		board.Description = description
	}

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}

	return board, nil
}

func (s *boardService) Delete(id uuid.UUID) error {
	_, err := s.boardRepo.FindByID(id)
	if err != nil {
		return ErrBoardNotFound
	}
	return s.boardRepo.Delete(id)
}

func (s *boardService) ListByProject(projectID uuid.UUID) ([]models.Board, error) {
	return s.boardRepo.ListByProject(projectID)
}
