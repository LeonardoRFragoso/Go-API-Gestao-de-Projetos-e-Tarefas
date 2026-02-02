package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService interface {
	Create(title, description string, priority models.TaskPriority, dueDate *time.Time, listID, creatorID uuid.UUID) (*models.Task, error)
	GetByID(id uuid.UUID) (*models.Task, error)
	GetByIDWithDetails(id uuid.UUID) (*models.Task, error)
	Update(id uuid.UUID, title, description string, priority models.TaskPriority, dueDate *time.Time) (*models.Task, error)
	Delete(id uuid.UUID) error
	ListByList(listID uuid.UUID) ([]models.Task, error)
	MoveTask(taskID, newListID uuid.UUID, newPosition int) error
	ReorderTasks(listID uuid.UUID, positions map[uuid.UUID]int) error
	AddAssignee(taskID, userID uuid.UUID) error
	RemoveAssignee(taskID, userID uuid.UUID) error
	AddLabel(taskID, labelID uuid.UUID) error
	RemoveLabel(taskID, labelID uuid.UUID) error
	Search(query string, projectID uuid.UUID, page, limit int) ([]models.Task, int64, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) Create(title, description string, priority models.TaskPriority, dueDate *time.Time, listID, creatorID uuid.UUID) (*models.Task, error) {
	maxPos, err := s.taskRepo.GetMaxPosition(listID)
	if err != nil {
		return nil, err
	}

	if priority == "" {
		priority = models.PriorityMedium
	}

	task := &models.Task{
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     dueDate,
		Position:    maxPos + 1,
		ListID:      listID,
		CreatorID:   creatorID,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return s.taskRepo.FindByIDWithDetails(task.ID)
}

func (s *taskService) GetByID(id uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (s *taskService) GetByIDWithDetails(id uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.FindByIDWithDetails(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (s *taskService) Update(id uuid.UUID, title, description string, priority models.TaskPriority, dueDate *time.Time) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if priority != "" {
		task.Priority = priority
	}
	if dueDate != nil {
		task.DueDate = dueDate
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return s.taskRepo.FindByIDWithDetails(task.ID)
}

func (s *taskService) Delete(id uuid.UUID) error {
	_, err := s.taskRepo.FindByID(id)
	if err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.Delete(id)
}

func (s *taskService) ListByList(listID uuid.UUID) ([]models.Task, error) {
	return s.taskRepo.ListByList(listID)
}

func (s *taskService) MoveTask(taskID, newListID uuid.UUID, newPosition int) error {
	_, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.MoveTask(taskID, newListID, newPosition)
}

func (s *taskService) ReorderTasks(listID uuid.UUID, positions map[uuid.UUID]int) error {
	return s.taskRepo.UpdatePositions(listID, positions)
}

func (s *taskService) AddAssignee(taskID, userID uuid.UUID) error {
	_, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.AddAssignee(taskID, userID)
}

func (s *taskService) RemoveAssignee(taskID, userID uuid.UUID) error {
	return s.taskRepo.RemoveAssignee(taskID, userID)
}

func (s *taskService) AddLabel(taskID, labelID uuid.UUID) error {
	_, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.AddLabel(taskID, labelID)
}

func (s *taskService) RemoveLabel(taskID, labelID uuid.UUID) error {
	return s.taskRepo.RemoveLabel(taskID, labelID)
}

func (s *taskService) Search(query string, projectID uuid.UUID, page, limit int) ([]models.Task, int64, error) {
	return s.taskRepo.Search(query, projectID, page, limit)
}
