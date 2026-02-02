package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	Color       string         `gorm:"default:'#8B5CF6'" json:"color"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Owner    User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members  []User    `gorm:"many2many:team_members;" json:"members,omitempty"`
	Projects []Project `gorm:"many2many:team_projects;" json:"projects,omitempty"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TeamRole string

const (
	TeamRoleLead   TeamRole = "lead"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleMember TeamRole = "member"
)

type TeamMember struct {
	TeamID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"team_id"`
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Role     TeamRole  `gorm:"type:varchar(20);default:'member'" json:"role"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`

	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type TeamProject struct {
	TeamID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"team_id"`
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey" json:"project_id"`
	AddedAt   time.Time `gorm:"autoCreateTime" json:"added_at"`

	Team    Team    `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

type TeamResponse struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Color       string         `json:"color"`
	OwnerID     uuid.UUID      `json:"owner_id"`
	Owner       *UserResponse  `json:"owner,omitempty"`
	MemberCount int            `json:"member_count"`
	ProjectCount int           `json:"project_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (t *Team) ToResponse() TeamResponse {
	resp := TeamResponse{
		ID:           t.ID,
		Name:         t.Name,
		Description:  t.Description,
		Color:        t.Color,
		OwnerID:      t.OwnerID,
		MemberCount:  len(t.Members),
		ProjectCount: len(t.Projects),
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
	if t.Owner.ID != uuid.Nil {
		ownerResp := t.Owner.ToResponse()
		resp.Owner = &ownerResp
	}
	return resp
}

type TeamMemberResponse struct {
	UserID   uuid.UUID     `json:"user_id"`
	User     *UserResponse `json:"user,omitempty"`
	Role     TeamRole      `json:"role"`
	JoinedAt time.Time     `json:"joined_at"`
}

type TeamDetailResponse struct {
	TeamResponse
	Members  []TeamMemberResponse `json:"members,omitempty"`
	Projects []ProjectResponse    `json:"projects,omitempty"`
}

func (t *Team) ToDetailResponse(members []TeamMember) TeamDetailResponse {
	resp := TeamDetailResponse{
		TeamResponse: t.ToResponse(),
		Members:      make([]TeamMemberResponse, 0, len(members)),
		Projects:     make([]ProjectResponse, 0, len(t.Projects)),
	}

	for _, m := range members {
		memberResp := TeamMemberResponse{
			UserID:   m.UserID,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		}
		if m.User.ID != uuid.Nil {
			userResp := m.User.ToResponse()
			memberResp.User = &userResp
		}
		resp.Members = append(resp.Members, memberResp)
	}

	for _, p := range t.Projects {
		resp.Projects = append(resp.Projects, p.ToResponse())
	}

	return resp
}
