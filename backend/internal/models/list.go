package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type List struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Position  int            `gorm:"not null;default:0" json:"position"`
	BoardID   uuid.UUID      `gorm:"type:uuid;not null" json:"board_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Board Board  `gorm:"foreignKey:BoardID" json:"board,omitempty"`
	Tasks []Task `gorm:"foreignKey:ListID" json:"tasks,omitempty"`
}

func (l *List) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type ListResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Position  int            `json:"position"`
	BoardID   uuid.UUID      `json:"board_id"`
	Tasks     []TaskResponse `json:"tasks,omitempty"`
	TaskCount int            `json:"task_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (l *List) ToResponse() ListResponse {
	tasks := make([]TaskResponse, 0, len(l.Tasks))
	for _, task := range l.Tasks {
		tasks = append(tasks, task.ToResponse())
	}
	return ListResponse{
		ID:        l.ID,
		Name:      l.Name,
		Position:  l.Position,
		BoardID:   l.BoardID,
		Tasks:     tasks,
		TaskCount: len(l.Tasks),
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}
