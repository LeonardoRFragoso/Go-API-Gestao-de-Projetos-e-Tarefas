package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Content   string         `gorm:"not null" json:"content"`
	TaskID    uuid.UUID      `gorm:"type:uuid;not null" json:"task_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Task Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type CommentResponse struct {
	ID        uuid.UUID     `json:"id"`
	Content   string        `json:"content"`
	TaskID    uuid.UUID     `json:"task_id"`
	UserID    uuid.UUID     `json:"user_id"`
	User      *UserResponse `json:"user,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (c *Comment) ToResponse() CommentResponse {
	resp := CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		TaskID:    c.TaskID,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.User.ID != uuid.Nil {
		userResp := c.User.ToResponse()
		resp.User = &userResp
	}
	return resp
}
