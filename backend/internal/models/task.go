package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description,omitempty"`
	Position    int            `gorm:"not null;default:0" json:"position"`
	Priority    TaskPriority   `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	ListID      uuid.UUID      `gorm:"type:uuid;not null" json:"list_id"`
	CreatorID   uuid.UUID      `gorm:"type:uuid;not null" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	List      List      `gorm:"foreignKey:ListID" json:"list,omitempty"`
	Creator   User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Assignees []User    `gorm:"many2many:task_assignees;" json:"assignees,omitempty"`
	Comments  []Comment `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
	Labels    []Label   `gorm:"many2many:task_labels;" json:"labels,omitempty"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TaskAssignee struct {
	TaskID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	AssignedAt time.Time `gorm:"autoCreateTime"`
}

type TaskResponse struct {
	ID           uuid.UUID      `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description,omitempty"`
	Position     int            `json:"position"`
	Priority     TaskPriority   `json:"priority"`
	DueDate      *time.Time     `json:"due_date,omitempty"`
	ListID       uuid.UUID      `json:"list_id"`
	CreatorID    uuid.UUID      `json:"creator_id"`
	Creator      *UserResponse  `json:"creator,omitempty"`
	Assignees    []UserResponse `json:"assignees,omitempty"`
	Labels       []LabelResponse `json:"labels,omitempty"`
	CommentCount int            `json:"comment_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (t *Task) ToResponse() TaskResponse {
	assignees := make([]UserResponse, 0, len(t.Assignees))
	for _, a := range t.Assignees {
		assignees = append(assignees, a.ToResponse())
	}

	labels := make([]LabelResponse, 0, len(t.Labels))
	for _, l := range t.Labels {
		labels = append(labels, l.ToResponse())
	}

	resp := TaskResponse{
		ID:           t.ID,
		Title:        t.Title,
		Description:  t.Description,
		Position:     t.Position,
		Priority:     t.Priority,
		DueDate:      t.DueDate,
		ListID:       t.ListID,
		CreatorID:    t.CreatorID,
		Assignees:    assignees,
		Labels:       labels,
		CommentCount: len(t.Comments),
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}

	if t.Creator.ID != uuid.Nil {
		creatorResp := t.Creator.ToResponse()
		resp.Creator = &creatorResp
	}

	return resp
}
