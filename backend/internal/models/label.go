package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Label struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Color     string         `gorm:"not null;default:'#6B7280'" json:"color"`
	BoardID   uuid.UUID      `gorm:"type:uuid;not null" json:"board_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Board Board  `gorm:"foreignKey:BoardID" json:"board,omitempty"`
	Tasks []Task `gorm:"many2many:task_labels;" json:"tasks,omitempty"`
}

func (l *Label) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type TaskLabel struct {
	TaskID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	LabelID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

type LabelResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Color   string    `json:"color"`
	BoardID uuid.UUID `json:"board_id"`
}

func (l *Label) ToResponse() LabelResponse {
	return LabelResponse{
		ID:      l.ID,
		Name:    l.Name,
		Color:   l.Color,
		BoardID: l.BoardID,
	}
}
