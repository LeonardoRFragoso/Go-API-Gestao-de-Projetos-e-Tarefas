package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTaskAssigned   NotificationType = "task_assigned"
	NotificationTaskUnassigned NotificationType = "task_unassigned"
	NotificationTaskUpdated    NotificationType = "task_updated"
	NotificationTaskComment    NotificationType = "task_comment"
	NotificationTaskDueSoon    NotificationType = "task_due_soon"
	NotificationTaskOverdue    NotificationType = "task_overdue"
	NotificationProjectInvite  NotificationType = "project_invite"
	NotificationTeamInvite     NotificationType = "team_invite"
	NotificationMention        NotificationType = "mention"
)

type Notification struct {
	ID        uuid.UUID        `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Title     string           `gorm:"not null" json:"title"`
	Message   string           `json:"message"`
	Read      bool             `gorm:"default:false" json:"read"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	
	// Reference to related entities
	TaskID    *uuid.UUID `gorm:"type:uuid" json:"task_id,omitempty"`
	ProjectID *uuid.UUID `gorm:"type:uuid" json:"project_id,omitempty"`
	TeamID    *uuid.UUID `gorm:"type:uuid" json:"team_id,omitempty"`
	ActorID   *uuid.UUID `gorm:"type:uuid" json:"actor_id,omitempty"`
	
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User  User `gorm:"foreignKey:UserID" json:"-"`
	Actor User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type NotificationResponse struct {
	ID        uuid.UUID        `json:"id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Read      bool             `json:"read"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	TaskID    *uuid.UUID       `json:"task_id,omitempty"`
	ProjectID *uuid.UUID       `json:"project_id,omitempty"`
	TeamID    *uuid.UUID       `json:"team_id,omitempty"`
	Actor     *UserResponse    `json:"actor,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

func (n *Notification) ToResponse() NotificationResponse {
	resp := NotificationResponse{
		ID:        n.ID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Read:      n.Read,
		ReadAt:    n.ReadAt,
		TaskID:    n.TaskID,
		ProjectID: n.ProjectID,
		TeamID:    n.TeamID,
		CreatedAt: n.CreatedAt,
	}
	if n.Actor.ID != uuid.Nil {
		actorResp := n.Actor.ToResponse()
		resp.Actor = &actorResp
	}
	return resp
}
