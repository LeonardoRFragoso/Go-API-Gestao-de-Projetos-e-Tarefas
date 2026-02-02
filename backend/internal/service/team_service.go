package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
)

var (
	ErrTeamNotFound      = errors.New("team not found")
	ErrNotTeamOwner      = errors.New("only team owner can perform this action")
	ErrAlreadyTeamMember = errors.New("user is already a team member")
	ErrNotTeamMember     = errors.New("user is not a team member")
	ErrCannotRemoveOwner = errors.New("cannot remove team owner")
)

type TeamService interface {
	Create(name, description, color string, ownerID uuid.UUID) (*models.Team, error)
	GetByID(id uuid.UUID) (*models.Team, error)
	GetByIDWithDetails(id uuid.UUID) (*models.Team, []models.TeamMember, error)
	Update(id uuid.UUID, name, description, color string) (*models.Team, error)
	Delete(id, userID uuid.UUID) error
	ListByUser(userID uuid.UUID, page, limit int) ([]models.Team, int64, error)
	
	// Member management
	AddMember(teamID, userID, actorID uuid.UUID, role models.TeamRole) error
	RemoveMember(teamID, userID, actorID uuid.UUID) error
	UpdateMemberRole(teamID, userID, actorID uuid.UUID, role models.TeamRole) error
	GetMembers(teamID uuid.UUID) ([]models.TeamMember, error)
	
	// Project management
	AddProject(teamID, projectID, actorID uuid.UUID) error
	RemoveProject(teamID, projectID, actorID uuid.UUID) error
	GetProjects(teamID uuid.UUID) ([]models.Project, error)
	
	// Permission check
	CanManageTeam(teamID, userID uuid.UUID) bool
}

type teamService struct {
	teamRepo repository.TeamRepository
	notificationService NotificationService
}

func NewTeamService(teamRepo repository.TeamRepository, notificationService NotificationService) TeamService {
	return &teamService{
		teamRepo:            teamRepo,
		notificationService: notificationService,
	}
}

func (s *teamService) Create(name, description, color string, ownerID uuid.UUID) (*models.Team, error) {
	if color == "" {
		color = "#8B5CF6"
	}

	team := &models.Team{
		Name:        name,
		Description: description,
		Color:       color,
		OwnerID:     ownerID,
	}

	if err := s.teamRepo.Create(team); err != nil {
		return nil, err
	}

	// Add owner as lead member
	if err := s.teamRepo.AddMember(team.ID, ownerID, models.TeamRoleLead); err != nil {
		return nil, err
	}

	return s.teamRepo.FindByID(team.ID)
}

func (s *teamService) GetByID(id uuid.UUID) (*models.Team, error) {
	team, err := s.teamRepo.FindByID(id)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}

func (s *teamService) GetByIDWithDetails(id uuid.UUID) (*models.Team, []models.TeamMember, error) {
	team, members, err := s.teamRepo.FindByIDWithDetails(id)
	if err != nil {
		return nil, nil, ErrTeamNotFound
	}
	return team, members, nil
}

func (s *teamService) Update(id uuid.UUID, name, description, color string) (*models.Team, error) {
	team, err := s.teamRepo.FindByID(id)
	if err != nil {
		return nil, ErrTeamNotFound
	}

	if name != "" {
		team.Name = name
	}
	if description != "" {
		team.Description = description
	}
	if color != "" {
		team.Color = color
	}

	if err := s.teamRepo.Update(team); err != nil {
		return nil, err
	}

	return s.teamRepo.FindByID(id)
}

func (s *teamService) Delete(id, userID uuid.UUID) error {
	team, err := s.teamRepo.FindByID(id)
	if err != nil {
		return ErrTeamNotFound
	}

	if team.OwnerID != userID {
		return ErrNotTeamOwner
	}

	return s.teamRepo.Delete(id)
}

func (s *teamService) ListByUser(userID uuid.UUID, page, limit int) ([]models.Team, int64, error) {
	return s.teamRepo.ListByUser(userID, page, limit)
}

func (s *teamService) AddMember(teamID, userID, actorID uuid.UUID, role models.TeamRole) error {
	if !s.CanManageTeam(teamID, actorID) {
		return ErrNotTeamOwner
	}

	if s.teamRepo.IsMember(teamID, userID) {
		return ErrAlreadyTeamMember
	}

	if err := s.teamRepo.AddMember(teamID, userID, role); err != nil {
		return err
	}

	// Send notification
	team, _ := s.teamRepo.FindByID(teamID)
	if team != nil && s.notificationService != nil {
		s.notificationService.NotifyTeamInvite(userID, team.ID, team.Name, actorID)
	}

	return nil
}

func (s *teamService) RemoveMember(teamID, userID, actorID uuid.UUID) error {
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return ErrTeamNotFound
	}

	if team.OwnerID == userID {
		return ErrCannotRemoveOwner
	}

	if !s.CanManageTeam(teamID, actorID) && actorID != userID {
		return ErrNotTeamOwner
	}

	if !s.teamRepo.IsMember(teamID, userID) {
		return ErrNotTeamMember
	}

	return s.teamRepo.RemoveMember(teamID, userID)
}

func (s *teamService) UpdateMemberRole(teamID, userID, actorID uuid.UUID, role models.TeamRole) error {
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return ErrTeamNotFound
	}

	if team.OwnerID == userID && role != models.TeamRoleLead {
		return ErrCannotRemoveOwner
	}

	if !s.CanManageTeam(teamID, actorID) {
		return ErrNotTeamOwner
	}

	return s.teamRepo.UpdateMemberRole(teamID, userID, role)
}

func (s *teamService) GetMembers(teamID uuid.UUID) ([]models.TeamMember, error) {
	return s.teamRepo.GetMembers(teamID)
}

func (s *teamService) AddProject(teamID, projectID, actorID uuid.UUID) error {
	if !s.CanManageTeam(teamID, actorID) {
		return ErrNotTeamOwner
	}

	return s.teamRepo.AddProject(teamID, projectID)
}

func (s *teamService) RemoveProject(teamID, projectID, actorID uuid.UUID) error {
	if !s.CanManageTeam(teamID, actorID) {
		return ErrNotTeamOwner
	}

	return s.teamRepo.RemoveProject(teamID, projectID)
}

func (s *teamService) GetProjects(teamID uuid.UUID) ([]models.Project, error) {
	return s.teamRepo.GetProjects(teamID)
}

func (s *teamService) CanManageTeam(teamID, userID uuid.UUID) bool {
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return false
	}

	if team.OwnerID == userID {
		return true
	}

	member, err := s.teamRepo.GetMember(teamID, userID)
	if err != nil {
		return false
	}

	return member.Role == models.TeamRoleLead || member.Role == models.TeamRoleAdmin
}
