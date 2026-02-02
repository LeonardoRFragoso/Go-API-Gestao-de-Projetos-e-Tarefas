package models

import (
	"github.com/google/uuid"
)

type Permission string

const (
	// Project permissions
	PermissionProjectView   Permission = "project:view"
	PermissionProjectEdit   Permission = "project:edit"
	PermissionProjectDelete Permission = "project:delete"
	PermissionProjectManageMembers Permission = "project:manage_members"

	// Board permissions
	PermissionBoardCreate Permission = "board:create"
	PermissionBoardEdit   Permission = "board:edit"
	PermissionBoardDelete Permission = "board:delete"

	// List permissions
	PermissionListCreate Permission = "list:create"
	PermissionListEdit   Permission = "list:edit"
	PermissionListDelete Permission = "list:delete"

	// Task permissions
	PermissionTaskCreate   Permission = "task:create"
	PermissionTaskEdit     Permission = "task:edit"
	PermissionTaskDelete   Permission = "task:delete"
	PermissionTaskAssign   Permission = "task:assign"
	PermissionTaskMove     Permission = "task:move"
	PermissionTaskComment  Permission = "task:comment"

	// Team permissions
	PermissionTeamView   Permission = "team:view"
	PermissionTeamEdit   Permission = "team:edit"
	PermissionTeamDelete Permission = "team:delete"
	PermissionTeamManageMembers Permission = "team:manage_members"
)

type ProjectRole string

const (
	ProjectRoleOwner  ProjectRole = "owner"
	ProjectRoleAdmin  ProjectRole = "admin"
	ProjectRoleMember ProjectRole = "member"
	ProjectRoleViewer ProjectRole = "viewer"
)

// RolePermissions defines which permissions each role has
var ProjectRolePermissions = map[ProjectRole][]Permission{
	ProjectRoleOwner: {
		PermissionProjectView, PermissionProjectEdit, PermissionProjectDelete, PermissionProjectManageMembers,
		PermissionBoardCreate, PermissionBoardEdit, PermissionBoardDelete,
		PermissionListCreate, PermissionListEdit, PermissionListDelete,
		PermissionTaskCreate, PermissionTaskEdit, PermissionTaskDelete, PermissionTaskAssign, PermissionTaskMove, PermissionTaskComment,
	},
	ProjectRoleAdmin: {
		PermissionProjectView, PermissionProjectEdit, PermissionProjectManageMembers,
		PermissionBoardCreate, PermissionBoardEdit, PermissionBoardDelete,
		PermissionListCreate, PermissionListEdit, PermissionListDelete,
		PermissionTaskCreate, PermissionTaskEdit, PermissionTaskDelete, PermissionTaskAssign, PermissionTaskMove, PermissionTaskComment,
	},
	ProjectRoleMember: {
		PermissionProjectView,
		PermissionBoardCreate, PermissionBoardEdit,
		PermissionListCreate, PermissionListEdit,
		PermissionTaskCreate, PermissionTaskEdit, PermissionTaskAssign, PermissionTaskMove, PermissionTaskComment,
	},
	ProjectRoleViewer: {
		PermissionProjectView,
		PermissionTaskComment,
	},
}

var TeamRolePermissions = map[TeamRole][]Permission{
	TeamRoleLead: {
		PermissionTeamView, PermissionTeamEdit, PermissionTeamDelete, PermissionTeamManageMembers,
	},
	TeamRoleAdmin: {
		PermissionTeamView, PermissionTeamEdit, PermissionTeamManageMembers,
	},
	TeamRoleMember: {
		PermissionTeamView,
	},
}

// HasPermission checks if a role has a specific permission
func HasProjectPermission(role ProjectRole, permission Permission) bool {
	permissions, exists := ProjectRolePermissions[role]
	if !exists {
		return false
	}
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func HasTeamPermission(role TeamRole, permission Permission) bool {
	permissions, exists := TeamRolePermissions[role]
	if !exists {
		return false
	}
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// Activity log for audit trail
type ActivityLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Action     string     `gorm:"not null" json:"action"`
	EntityType string     `gorm:"not null" json:"entity_type"`
	EntityID   uuid.UUID  `gorm:"type:uuid" json:"entity_id"`
	Details    string     `gorm:"type:text" json:"details,omitempty"`
	IPAddress  string     `json:"ip_address,omitempty"`
	CreatedAt  string     `json:"created_at"`
}
