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
	Move(id uuid.UUID, targetListID uuid.UUID, position int) (*models.Task, error)
	MoveTask(taskID, newListID uuid.UUID, newPosition int) error
	Reorder(listID uuid.UUID, taskIDs []uuid.UUID) error
	ReorderTasks(listID uuid.UUID, positions map[uuid.UUID]int) error
	Search(query string, projectID uuid.UUID, page, limit int) ([]*models.Task, int64, error)
	AddAssignee(taskID, userID uuid.UUID) error
	RemoveAssignee(taskID, userID uuid.UUID) error
	AddLabel(taskID, labelID uuid.UUID) error
	RemoveLabel(taskID, labelID uuid.UUID) error
}

type taskServiceImpl struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskServiceImpl{taskRepo: taskRepo}
}

func (s *taskServiceImpl) Create(title, description string, priority models.TaskPriority, dueDate *time.Time, listID, creatorID uuid.UUID) (*models.Task, error) {
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

func (s *taskServiceImpl) GetByID(id uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (s *taskServiceImpl) GetByIDWithDetails(id uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.FindByIDWithDetails(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (s *taskServiceImpl) Update(id uuid.UUID, title, description string, priority models.TaskPriority, dueDate *time.Time) (*models.Task, error) {
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

func (s *taskServiceImpl) Delete(id uuid.UUID) error {
	_, err := s.taskRepo.FindByID(id)
	if err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.Delete(id)
}

func (s *taskServiceImpl) ListByList(listID uuid.UUID) ([]models.Task, error) {
	return s.taskRepo.ListByList(listID)
}

func (s *taskServiceImpl) Move(id uuid.UUID, targetListID uuid.UUID, position int) (*models.Task, error) {
	_, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	err = s.taskRepo.MoveTask(id, targetListID, position)
	if err != nil {
		return nil, err
	}
	return s.taskRepo.FindByIDWithDetails(id)
}

func (s *taskServiceImpl) MoveTask(taskID, newListID uuid.UUID, newPosition int) error {
	_, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.MoveTask(taskID, newListID, newPosition)
}

func (s *taskServiceImpl) Reorder(listID uuid.UUID, taskIDs []uuid.UUID) error {
	positions := make(map[uuid.UUID]int)
	for i, id := range taskIDs {
		positions[id] = i
	}
	return s.taskRepo.UpdatePositions(listID, positions)
}

func (s *taskServiceImpl) ReorderTasks(listID uuid.UUID, positions map[uuid.UUID]int) error {
	return s.taskRepo.UpdatePositions(listID, positions)
}

func (s *taskServiceImpl) AddAssignee(taskID, userID uuid.UUID) error {
	_, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.AddAssignee(taskID, userID)
}

func (s *taskServiceImpl) RemoveAssignee(taskID, userID uuid.UUID) error {
	return s.taskRepo.RemoveAssignee(taskID, userID)
}

func (s *taskServiceImpl) AddLabel(taskID, labelID uuid.UUID) error {
	_, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.AddLabel(taskID, labelID)
}

func (s *taskServiceImpl) RemoveLabel(taskID, labelID uuid.UUID) error {
	return s.taskRepo.RemoveLabel(taskID, labelID)
}

func (s *taskServiceImpl) Search(query string, projectID uuid.UUID, page, limit int) ([]*models.Task, int64, error) {
	tasks, total, err := s.taskRepo.Search(query, projectID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*models.Task, len(tasks))
	for i := range tasks {
		result[i] = &tasks[i]
	}
	return result, total, nil
}

func (s *taskServiceImpl) GetTaskStatsByUser(userID uuid.UUID) (map[string]int, error) {
	return s.taskRepo.GetTaskStatsByUser(userID)
}

func (s *taskServiceImpl) GetWeeklyCompletedTasks(userID uuid.UUID) ([]int, error) {
	return s.taskRepo.GetWeeklyCompletedTasks(userID)
}
