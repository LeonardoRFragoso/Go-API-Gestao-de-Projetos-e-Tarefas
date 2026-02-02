package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrNotCommentOwner = errors.New("user is not comment owner")
)

type CommentService interface {
	Create(content string, taskID, userID uuid.UUID) (*models.Comment, error)
	GetByID(id uuid.UUID) (*models.Comment, error)
	Update(id, userID uuid.UUID, content string) (*models.Comment, error)
	Delete(id, userID uuid.UUID) error
	ListByTask(taskID uuid.UUID, page, limit int) ([]models.Comment, int64, error)
}

type commentService struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) CommentService {
	return &commentService{commentRepo: commentRepo}
}

func (s *commentService) Create(content string, taskID, userID uuid.UUID) (*models.Comment, error) {
	comment := &models.Comment{
		Content: content,
		TaskID:  taskID,
		UserID:  userID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return s.commentRepo.FindByID(comment.ID)
}

func (s *commentService) GetByID(id uuid.UUID) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return nil, ErrCommentNotFound
	}
	return comment, nil
}

func (s *commentService) Update(id, userID uuid.UUID, content string) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return nil, ErrCommentNotFound
	}

	if comment.UserID != userID {
		return nil, ErrNotCommentOwner
	}

	comment.Content = content
	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) Delete(id, userID uuid.UUID) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return ErrCommentNotFound
	}

	if comment.UserID != userID {
		return ErrNotCommentOwner
	}

	return s.commentRepo.Delete(id)
}

func (s *commentService) ListByTask(taskID uuid.UUID, page, limit int) ([]models.Comment, int64, error) {
	return s.commentRepo.ListByTask(taskID, page, limit)
}
