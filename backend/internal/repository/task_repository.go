package repository

import (
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	FindByID(id uuid.UUID) (*models.Task, error)
	FindByIDWithDetails(id uuid.UUID) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id uuid.UUID) error
	ListByList(listID uuid.UUID) ([]models.Task, error)
	GetMaxPosition(listID uuid.UUID) (int, error)
	UpdatePositions(listID uuid.UUID, positions map[uuid.UUID]int) error
	MoveTask(taskID, newListID uuid.UUID, newPosition int) error
	AddAssignee(taskID, userID uuid.UUID) error
	RemoveAssignee(taskID, userID uuid.UUID) error
	AddLabel(taskID, labelID uuid.UUID) error
	RemoveLabel(taskID, labelID uuid.UUID) error
	Search(query string, projectID uuid.UUID, page, limit int) ([]models.Task, int64, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uuid.UUID) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindByIDWithDetails(id uuid.UUID) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("Creator").
		Preload("Assignees").
		Preload("Labels").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Comments.User").
		First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Task{}, "id = ?", id).Error
}

func (r *taskRepository) ListByList(listID uuid.UUID) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Preload("Assignees").Preload("Labels").
		Where("list_id = ?", listID).
		Order("position ASC").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetMaxPosition(listID uuid.UUID) (int, error) {
	var maxPos int
	err := r.db.Model(&models.Task{}).
		Where("list_id = ?", listID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPos).Error
	return maxPos, err
}

func (r *taskRepository) UpdatePositions(listID uuid.UUID, positions map[uuid.UUID]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for taskID, pos := range positions {
			if err := tx.Model(&models.Task{}).
				Where("id = ? AND list_id = ?", taskID, listID).
				Update("position", pos).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *taskRepository) MoveTask(taskID, newListID uuid.UUID, newPosition int) error {
	return r.db.Model(&models.Task{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"list_id":  newListID,
			"position": newPosition,
		}).Error
}

func (r *taskRepository) AddAssignee(taskID, userID uuid.UUID) error {
	assignee := models.TaskAssignee{
		TaskID: taskID,
		UserID: userID,
	}
	return r.db.Create(&assignee).Error
}

func (r *taskRepository) RemoveAssignee(taskID, userID uuid.UUID) error {
	return r.db.Delete(&models.TaskAssignee{}, "task_id = ? AND user_id = ?", taskID, userID).Error
}

func (r *taskRepository) AddLabel(taskID, labelID uuid.UUID) error {
	label := models.TaskLabel{
		TaskID:  taskID,
		LabelID: labelID,
	}
	return r.db.Create(&label).Error
}

func (r *taskRepository) RemoveLabel(taskID, labelID uuid.UUID) error {
	return r.db.Delete(&models.TaskLabel{}, "task_id = ? AND label_id = ?", taskID, labelID).Error
}

func (r *taskRepository) Search(query string, projectID uuid.UUID, page, limit int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	offset := (page - 1) * limit

	baseQuery := r.db.Model(&models.Task{}).
		Joins("JOIN lists ON lists.id = tasks.list_id").
		Joins("JOIN boards ON boards.id = lists.board_id").
		Where("boards.project_id = ?", projectID).
		Where("tasks.title ILIKE ? OR tasks.description ILIKE ?", "%"+query+"%", "%"+query+"%")

	err := baseQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Assignees").Preload("Labels").
		Joins("JOIN lists ON lists.id = tasks.list_id").
		Joins("JOIN boards ON boards.id = lists.board_id").
		Where("boards.project_id = ?", projectID).
		Where("tasks.title ILIKE ? OR tasks.description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Offset(offset).Limit(limit).
		Find(&tasks).Error

	return tasks, total, err
}
