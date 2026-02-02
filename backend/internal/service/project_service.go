package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrNotProjectOwner = errors.New("user is not project owner")
	ErrNotProjectMember = errors.New("user is not project member")
)

type ProjectService interface {
	Create(name, description, color string, ownerID uuid.UUID) (*models.Project, error)
	GetByID(id uuid.UUID) (*models.Project, error)
	GetByIDWithDetails(id uuid.UUID) (*models.Project, error)
	Update(id uuid.UUID, userID uuid.UUID, name, description, color string) (*models.Project, error)
	Delete(id, userID uuid.UUID) error
	ListByUser(userID uuid.UUID, page, limit int) ([]models.Project, int64, error)
	AddMember(projectID, userID, memberID uuid.UUID, role string) error
	RemoveMember(projectID, userID, memberID uuid.UUID) error
	GetMembers(projectID uuid.UUID) ([]models.User, error)
	CanAccess(projectID, userID uuid.UUID) (bool, error)
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

func (s *projectService) Create(name, description, color string, ownerID uuid.UUID) (*models.Project, error) {
	project := &models.Project{
		Name:        name,
		Description: description,
		Color:       color,
		OwnerID:     ownerID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	if err := s.projectRepo.AddMember(project.ID, ownerID, "owner"); err != nil {
		return nil, err
	}

	return s.projectRepo.FindByIDWithDetails(project.ID)
}

func (s *projectService) GetByID(id uuid.UUID) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

func (s *projectService) GetByIDWithDetails(id uuid.UUID) (*models.Project, error) {
	project, err := s.projectRepo.FindByIDWithDetails(id)
	if err != nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

func (s *projectService) Update(id uuid.UUID, userID uuid.UUID, name, description, color string) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, ErrProjectNotFound
	}

	if project.OwnerID != userID {
		return nil, ErrNotProjectOwner
	}

	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}
	if color != "" {
		project.Color = color
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return s.projectRepo.FindByIDWithDetails(project.ID)
}

func (s *projectService) Delete(id, userID uuid.UUID) error {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return ErrProjectNotFound
	}

	if project.OwnerID != userID {
		return ErrNotProjectOwner
	}

	return s.projectRepo.Delete(id)
}

func (s *projectService) ListByUser(userID uuid.UUID, page, limit int) ([]models.Project, int64, error) {
	return s.projectRepo.ListByUser(userID, page, limit)
}

func (s *projectService) AddMember(projectID, userID, memberID uuid.UUID, role string) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return ErrProjectNotFound
	}

	if project.OwnerID != userID {
		return ErrNotProjectOwner
	}

	return s.projectRepo.AddMember(projectID, memberID, role)
}

func (s *projectService) RemoveMember(projectID, userID, memberID uuid.UUID) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return ErrProjectNotFound
	}

	if project.OwnerID != userID {
		return ErrNotProjectOwner
	}

	if project.OwnerID == memberID {
		return errors.New("cannot remove project owner")
	}

	return s.projectRepo.RemoveMember(projectID, memberID)
}

func (s *projectService) GetMembers(projectID uuid.UUID) ([]models.User, error) {
	return s.projectRepo.GetMembers(projectID)
}

func (s *projectService) CanAccess(projectID, userID uuid.UUID) (bool, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return false, ErrProjectNotFound
	}

	if project.OwnerID == userID {
		return true, nil
	}

	return s.projectRepo.IsMember(projectID, userID)
}
