package types

import (
	"slices"
)

// Action represents possible actions on resources
type Action string

const (
	ActionCreate        Action = "create"
	ActionDelete        Action = "delete"
	ActionDisable       Action = "disable"
	ActionEnable        Action = "enable"
	ActionManage        Action = "manage" // Implies full access
	ActionRead          Action = "read"
	ActionResetPassword Action = "reset_password"
	ActionUpdate        Action = "update"
	ActionVerify        Action = "verify"
)

// String returns the string representation of an action
func (a Action) String() string {
	return string(a)
}

// Resource represents protected resources
type Resource string

const (
	ResourceEntity    Resource = "entity"
	ResourcePayment   Resource = "payment"
	ResourceSession   Resource = "session"
	ResourceTwoFactor Resource = "two_factor"
	ResourceUser      Resource = "user"
)

// String returns the string representation of a resource
func (r Resource) String() string {
	return string(r)
}

// Role represents predefined roles in the system
type Role string

const (
	RoleNone   Role = ""        // No access
	RoleOwner  Role = "owner"   // Full access to everything
	RoleAPIKey Role = "api-key" // Secret key access
	RoleAdmin  Role = "admin"   // Full access except critical operations
	RoleViewer Role = "viewer"  // Read-only access
)

// String returns the string representation of a role
func (r Role) String() string {
	return string(r)
}

func (r Role) HasPermission(resource Resource, action Action) bool {
	resourcePerms, exists := RolePermissions[r]
	if !exists {
		return false
	}

	actions, exists := resourcePerms[resource]
	if !exists {
		return false
	}

	if slices.Contains(actions, ActionManage) {
		return true
	}

	return slices.Contains(actions, action)
}

// GetPermissions returns all permissions for a role
func (r Role) GetPermissions() []Permission {
	resourcePerms, exists := RolePermissions[r]
	if !exists {
		return nil
	}

	var permissions []Permission
	for resource, actions := range resourcePerms {
		for _, action := range actions {
			permissions = append(permissions, Permission{
				Resource: resource,
				Action:   action,
			})
		}
	}

	return permissions
}

// Permission represents a permission to perform an action on a resource
type Permission struct {
	Resource Resource `json:"resource"`
	Action   Action   `json:"action"`
}

// RolePermissions defines the permissions for each role
var RolePermissions = map[Role]map[Resource][]Action{
	RoleOwner: {
		// Full access to everything
		ResourceEntity:  {ActionManage},
		ResourceUser:    {ActionManage},
		ResourcePayment: {ActionManage},
	},
	RoleAdmin: {
		// Full access except critical operations
		ResourceEntity:  {ActionRead, ActionUpdate},
		ResourceUser:    {ActionManage},
		ResourcePayment: {ActionManage},
	},
	RoleViewer: {
		// Read-only access
		ResourceEntity:  {ActionRead},
		ResourceUser:    {ActionManage},
		ResourcePayment: {ActionRead},
	},

	// API Key access for most services
	RoleAPIKey: {
		ResourceEntity:  {ActionRead},
		ResourceUser:    {ActionRead},
		ResourcePayment: {ActionManage},
	},
}

// IsValidRole checks if a role is valid
func IsValidRole(role Role) bool {
	_, exists := RolePermissions[role]

	return exists
}

// GetAvailableRoles returns all available roles
func GetAvailableRoles() []Role {
	roles := make([]Role, 0, len(RolePermissions))
	for role := range RolePermissions {
		roles = append(roles, role)
	}
	return roles
}
