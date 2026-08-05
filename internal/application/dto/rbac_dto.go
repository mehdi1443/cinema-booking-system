package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type RoleResponse struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Title       string                `json:"title"`
	CreatedAt   *time.Time            `json:"created_at,omitempty"`
	UpdatedAt   *time.Time            `json:"updated_at,omitempty"`
	DeletedAt   *time.Time            `json:"deleted_at,omitempty"`
	Permissions []*PermissionResponse `json:"permissions,omitempty"`
	Users       []*UserResponse       `json:"users,omitempty"`
}

type UpdateUserRolesRequest struct {
	RoleIDs []int64 `json:"role_ids" binding:"required"`
}

type AssignPermissionsToRoleRequest struct {
	PermissionIDs []*int64 `json:"permission_ids" binding:"required"`
}

type DischargePermissionsOfRoleRequest struct {
	PermissionIDs []*int64 `json:"permission_ids" binding:"required"`
}

func ToRoleResponse(role *entity.Role) *RoleResponse {
	if role == nil {
		return nil
	}

	res := &RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		Title:     role.Title,
		CreatedAt: helper.TimeToPtr(role.CreatedAt),
		UpdatedAt: helper.TimeToPtr(role.UpdatedAt),
		DeletedAt: helper.TimeToPtr(role.DeletedAt.Time),
	}

	//HACK : Check it...
	// if !role.DeletedAt.IsZero() {
	// 	res.DeletedAt = helper.TimeToPtr(role.DeletedAt)
	// }

	// Only populate permissions if relationship was preloaded
	if len(role.Permissions) > 0 {
		res.Permissions = ToPermissionResponseSlice(role.Permissions)
	}
	if len(role.Users) > 0 {
		res.Users = ToUserResponseSlice(role.Users, &HideFieldUserOption{WalletBalance: false})
	}

	return res
}

func ToRoleResponseSlice(roles []*entity.Role) []*RoleResponse {
	if len(roles) == 0 {
		return []*RoleResponse{}
	}

	responses := make([]*RoleResponse, 0, len(roles))
	for ind := range roles {
		if resp := ToRoleResponse(roles[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}
