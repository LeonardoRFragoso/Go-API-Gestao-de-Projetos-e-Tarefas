package repository

import (
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uuid.UUID) (*models.Project, error)
	FindByIDWithDetails(id uuid.UUID) (*models.Project, error)
	Update(project *models.Project) error
	Delete(id uuid.UUID) error
	ListByUser(userID uuid.UUID, page, limit int) ([]models.Project, int64, error)
	AddMember(projectID, userID uuid.UUID, role string) error
	RemoveMember(projectID, userID uuid.UUID) error
	GetMembers(projectID uuid.UUID) ([]models.User, error)
	IsMember(projectID, userID uuid.UUID) (bool, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Owner").First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindByIDWithDetails(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Owner").Preload("Members").Preload("Boards").First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}

func (r *projectRepository) ListByUser(userID uuid.UUID, page, limit int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	offset := (page - 1) * limit

	subQuery := r.db.Table("project_members").Select("project_id").Where("user_id = ?", userID)

	countQuery := r.db.Model(&models.Project{}).Where("owner_id = ? OR id IN (?)", userID, subQuery)
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Owner").Preload("Members").Preload("Boards").
		Where("owner_id = ? OR id IN (?)", userID, subQuery).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) AddMember(projectID, userID uuid.UUID, role string) error {
	member := models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	return r.db.Create(&member).Error
}

func (r *projectRepository) RemoveMember(projectID, userID uuid.UUID) error {
	return r.db.Delete(&models.ProjectMember{}, "project_id = ? AND user_id = ?", projectID, userID).Error
}

func (r *projectRepository) GetMembers(projectID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Joins("JOIN project_members ON project_members.user_id = users.id").
		Where("project_members.project_id = ?", projectID).
		Find(&users).Error
	return users, err
}

func (r *projectRepository) IsMember(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}
