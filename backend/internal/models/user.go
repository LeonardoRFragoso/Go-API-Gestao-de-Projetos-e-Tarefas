package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	Avatar    string         `json:"avatar,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Projects       []Project       `gorm:"many2many:project_members;" json:"projects,omitempty"`
	OwnedProjects  []Project       `gorm:"foreignKey:OwnerID" json:"owned_projects,omitempty"`
	AssignedTasks  []Task          `gorm:"many2many:task_assignees;" json:"assigned_tasks,omitempty"`
	Comments       []Comment       `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	RefreshTokens  []RefreshToken  `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
	}
}
