package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
)

type NotificationService interface {
	Create(userID uuid.UUID, notifType models.NotificationType, title, message string, taskID, projectID, teamID, actorID *uuid.UUID) error
	GetByID(id uuid.UUID) (*models.Notification, error)
	ListByUser(userID uuid.UUID, page, limit int, unreadOnly bool) ([]models.Notification, int64, error)
	MarkAsRead(id, userID uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	Delete(id, userID uuid.UUID) error
	DeleteAll(userID uuid.UUID) error
	GetUnreadCount(userID uuid.UUID) (int64, error)

	// Notification helpers
	NotifyTaskAssigned(userID, taskID uuid.UUID, taskTitle string, actorID uuid.UUID) error
	NotifyTaskUnassigned(userID, taskID uuid.UUID, taskTitle string, actorID uuid.UUID) error
	NotifyTaskUpdated(userID, taskID uuid.UUID, taskTitle string, actorID uuid.UUID) error
	NotifyTaskComment(userID, taskID uuid.UUID, taskTitle, comment string, actorID uuid.UUID) error
	NotifyProjectInvite(userID, projectID uuid.UUID, projectName string, actorID uuid.UUID) error
	NotifyTeamInvite(userID, teamID uuid.UUID, teamName string, actorID uuid.UUID) error
	NotifyMention(userID, taskID uuid.UUID, taskTitle string, actorID uuid.UUID) error
}

type notificationService struct {
	notifRepo repository.NotificationRepository
}

func NewNotificationService(notifRepo repository.NotificationRepository) NotificationService {
	return &notificationService{notifRepo: notifRepo}
}

func (s *notificationService) Create(userID uuid.UUID, notifType models.NotificationType, title, message string, taskID, projectID, teamID, actorID *uuid.UUID) error {
	notification := &models.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		TaskID:    taskID,
		ProjectID: projectID,
		TeamID:    teamID,
		ActorID:   actorID,
	}
	return s.notifRepo.Create(notification)
}

func (s *notificationService) GetByID(id uuid.UUID) (*models.Notification, error) {
	return s.notifRepo.FindByID(id)
}

func (s *notificationService) ListByUser(userID uuid.UUID, page, limit int, unreadOnly bool) ([]models.Notification, int64, error) {
	return s.notifRepo.ListByUser(userID, page, limit, unreadOnly)
}

func (s *notificationService) MarkAsRead(id, userID uuid.UUID) error {
	notification, err := s.notifRepo.FindByID(id)
	if err != nil {
		return err
	}
	if notification.UserID != userID {
		return fmt.Errorf("notification not found")
	}
	return s.notifRepo.MarkAsRead(id)
}

func (s *notificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notifRepo.MarkAllAsRead(userID)
}

func (s *notificationService) Delete(id, userID uuid.UUID) error {
	notification, err := s.notifRepo.FindByID(id)
	if err != nil {
		return err
	}
	if notification.UserID != userID {
		return fmt.Errorf("notification not found")
	}
	return s.notifRepo.Delete(id)
}

func (s *notificationService) DeleteAll(userID uuid.UUID) error {
	return s.notifRepo.DeleteAllByUser(userID)
}

func (s *notificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.notifRepo.GetUnreadCount(userID)
}

// Helper methods for common notifications

func (s *notificationService) NotifyTaskAssigned(userID, taskID uuid.UUID, taskTitle string, actorID uuid.UUID) error {
	return s.Create(
		userID,
		models.NotificationTaskAssigned,
		"Tarefa atribuída",
		fmt.Sprintf("Você foi atribuído à tarefa: %s", taskTitle),
		&taskID, nil, nil, &actorID,
	)
}

func (s *notificationService) NotifyTaskUnassigned(userID, taskID uuid.UUID, taskTitle string, actorID uuid.UUID) error {
	return s.Create(
		userID,
		models.NotificationTaskUnassigned,
		"Tarefa removida",
		fmt.Sprintf("Você foi removido da tarefa: %s", taskTitle),
		&taskID, nil, nil, &actorID,
	)
}

func (s *notificationService) NotifyTaskUpdated(userID, taskID uuid.UUID, taskTitle string, actorID uuid.UUID) error {
	return s.Create(
		userID,
		models.NotificationTaskUpdated,
		"Tarefa atualizada",
		fmt.Sprintf("A tarefa '%s' foi atualizada", taskTitle),
		&taskID, nil, nil, &actorID,
	)
}

func (s *notificationService) NotifyTaskComment(userID, taskID uuid.UUID, taskTitle, comment string, actorID uuid.UUID) error {
	return s.Create(
		userID,
		models.NotificationTaskComment,
		"Novo comentário",
		fmt.Sprintf("Novo comentário na tarefa '%s': %s", taskTitle, truncate(comment, 100)),
		&taskID, nil, nil, &actorID,
	)
}

func (s *notificationService) NotifyProjectInvite(userID, projectID uuid.UUID, projectName string, actorID uuid.UUID) error {
	return s.Create(
		userID,
		models.NotificationProjectInvite,
		"Convite para projeto",
		fmt.Sprintf("Você foi convidado para o projeto: %s", projectName),
		nil, &projectID, nil, &actorID,
	)
}

func (s *notificationService) NotifyTeamInvite(userID, teamID uuid.UUID, teamName string, actorID uuid.UUID) error {
	return s.Create(
		userID,
		models.NotificationTeamInvite,
		"Convite para equipe",
		fmt.Sprintf("Você foi convidado para a equipe: %s", teamName),
		nil, nil, &teamID, &actorID,
	)
}

func (s *notificationService) NotifyMention(userID, taskID uuid.UUID, taskTitle string, actorID uuid.UUID) error {
	return s.Create(
		userID,
		models.NotificationMention,
		"Você foi mencionado",
		fmt.Sprintf("Você foi mencionado na tarefa: %s", taskTitle),
		&taskID, nil, nil, &actorID,
	)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
