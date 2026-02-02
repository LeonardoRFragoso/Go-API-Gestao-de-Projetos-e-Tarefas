package repository

import (
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(team *models.Team) error
	FindByID(id uuid.UUID) (*models.Team, error)
	FindByIDWithDetails(id uuid.UUID) (*models.Team, []models.TeamMember, error)
	Update(team *models.Team) error
	Delete(id uuid.UUID) error
	ListByUser(userID uuid.UUID, page, limit int) ([]models.Team, int64, error)
	
	// Member management
	AddMember(teamID, userID uuid.UUID, role models.TeamRole) error
	RemoveMember(teamID, userID uuid.UUID) error
	UpdateMemberRole(teamID, userID uuid.UUID, role models.TeamRole) error
	GetMember(teamID, userID uuid.UUID) (*models.TeamMember, error)
	GetMembers(teamID uuid.UUID) ([]models.TeamMember, error)
	IsMember(teamID, userID uuid.UUID) bool
	
	// Project management
	AddProject(teamID, projectID uuid.UUID) error
	RemoveProject(teamID, projectID uuid.UUID) error
	GetProjects(teamID uuid.UUID) ([]models.Project, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepository) FindByID(id uuid.UUID) (*models.Team, error) {
	var team models.Team
	err := r.db.Preload("Owner").First(&team, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) FindByIDWithDetails(id uuid.UUID) (*models.Team, []models.TeamMember, error) {
	var team models.Team
	err := r.db.Preload("Owner").Preload("Members").Preload("Projects").First(&team, "id = ?", id).Error
	if err != nil {
		return nil, nil, err
	}

	var members []models.TeamMember
	err = r.db.Preload("User").Where("team_id = ?", id).Find(&members).Error
	if err != nil {
		return nil, nil, err
	}

	return &team, members, nil
}

func (r *teamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

func (r *teamRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Team{}, "id = ?", id).Error
}

func (r *teamRepository) ListByUser(userID uuid.UUID, page, limit int) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	offset := (page - 1) * limit

	// Get teams where user is owner or member
	subQuery := r.db.Table("team_members").Select("team_id").Where("user_id = ?", userID)
	
	query := r.db.Model(&models.Team{}).
		Where("owner_id = ? OR id IN (?)", userID, subQuery)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Owner").Preload("Members").Preload("Projects").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&teams).Error

	return teams, total, err
}

func (r *teamRepository) AddMember(teamID, userID uuid.UUID, role models.TeamRole) error {
	member := models.TeamMember{
		TeamID: teamID,
		UserID: userID,
		Role:   role,
	}
	return r.db.Create(&member).Error
}

func (r *teamRepository) RemoveMember(teamID, userID uuid.UUID) error {
	return r.db.Delete(&models.TeamMember{}, "team_id = ? AND user_id = ?", teamID, userID).Error
}

func (r *teamRepository) UpdateMemberRole(teamID, userID uuid.UUID, role models.TeamRole) error {
	return r.db.Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Update("role", role).Error
}

func (r *teamRepository) GetMember(teamID, userID uuid.UUID) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.Preload("User").First(&member, "team_id = ? AND user_id = ?", teamID, userID).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *teamRepository) GetMembers(teamID uuid.UUID) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.db.Preload("User").Where("team_id = ?", teamID).Find(&members).Error
	return members, err
}

func (r *teamRepository) IsMember(teamID, userID uuid.UUID) bool {
	var count int64
	r.db.Model(&models.TeamMember{}).Where("team_id = ? AND user_id = ?", teamID, userID).Count(&count)
	return count > 0
}

func (r *teamRepository) AddProject(teamID, projectID uuid.UUID) error {
	tp := models.TeamProject{
		TeamID:    teamID,
		ProjectID: projectID,
	}
	return r.db.Create(&tp).Error
}

func (r *teamRepository) RemoveProject(teamID, projectID uuid.UUID) error {
	return r.db.Delete(&models.TeamProject{}, "team_id = ? AND project_id = ?", teamID, projectID).Error
}

func (r *teamRepository) GetProjects(teamID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Joins("JOIN team_projects ON team_projects.project_id = projects.id").
		Where("team_projects.team_id = ?", teamID).
		Preload("Owner").
		Find(&projects).Error
	return projects, err
}
