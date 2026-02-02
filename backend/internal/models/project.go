package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	Color       string         `gorm:"default:'#3B82F6'" json:"color"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Owner   User    `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members []User  `gorm:"many2many:project_members;" json:"members,omitempty"`
	Boards  []Board `gorm:"foreignKey:ProjectID" json:"boards,omitempty"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type ProjectMember struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Role      string    `gorm:"default:'member'"` // owner, admin, member
	JoinedAt  time.Time `gorm:"autoCreateTime"`
}

type ProjectResponse struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Color       string         `json:"color"`
	OwnerID     uuid.UUID      `json:"owner_id"`
	Owner       *UserResponse  `json:"owner,omitempty"`
	MemberCount int            `json:"member_count,omitempty"`
	BoardCount  int            `json:"board_count,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (p *Project) ToResponse() ProjectResponse {
	resp := ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Color:       p.Color,
		OwnerID:     p.OwnerID,
		MemberCount: len(p.Members),
		BoardCount:  len(p.Boards),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Owner.ID != uuid.Nil {
		ownerResp := p.Owner.ToResponse()
		resp.Owner = &ownerResp
	}
	return resp
}
