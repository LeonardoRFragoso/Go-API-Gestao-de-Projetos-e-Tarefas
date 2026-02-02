package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;not null" json:"project_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Project Project  `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Lists   []List   `gorm:"foreignKey:BoardID" json:"lists,omitempty"`
}

func (b *Board) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BoardResponse struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	ProjectID   uuid.UUID      `json:"project_id"`
	ListCount   int            `json:"list_count,omitempty"`
	TaskCount   int            `json:"task_count,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (b *Board) ToResponse() BoardResponse {
	taskCount := 0
	for _, list := range b.Lists {
		taskCount += len(list.Tasks)
	}
	return BoardResponse{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		ProjectID:   b.ProjectID,
		ListCount:   len(b.Lists),
		TaskCount:   taskCount,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
