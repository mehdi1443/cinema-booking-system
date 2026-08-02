package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

// [ ] : Im not Sure but i think i declare PermissionID on its enum_number instead of DB ID ! Check it!!! same it for Category (confused it is string or int)
type PermissionResponse struct {
	ID          int64                   `json:"id"`
	Name        string                  `json:"name"`
	Title       string                  `json:"title"`
	Category    enum.PermissionCategory `json:"category"`
	Description string                  `json:"description,omitempty"`
	CreatedAt   *time.Time              `json:"created_at,omitempty"`
	UpdatedAt   *time.Time              `json:"updated_at,omitempty"`
	DeletedAt   *time.Time              `json:"deleted_at,omitempty"`
	Roles       []*RoleResponse         `json:"roles,omitempty"`
}

func ToPermissionResponse(perm *entity.Permission) *PermissionResponse {
	if perm == nil {
		return nil
	}

	res := &PermissionResponse{
		ID:          perm.ID,
		Name:        perm.Name,
		Title:       perm.Title,
		Category:    perm.Category,
		Description: perm.Description,
		CreatedAt:   helper.TimeToPtr(perm.CreatedAt),
		UpdatedAt:   helper.TimeToPtr(perm.UpdatedAt),
		DeletedAt:   helper.TimeToPtr(perm.DeletedAt.Time),
	}

	if len(res.Roles) > 0 {
		res.Roles = ToRoleResponseSlice(perm.Roles)
	}

	return res
}

func ToPermissionResponseSlice(perms []*entity.Permission) []*PermissionResponse {
	if len(perms) == 0 {
		return []*PermissionResponse{}
	}

	responses := make([]*PermissionResponse, 0, len(perms))
	for ind := range perms {
		if resp := ToPermissionResponse(perms[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}
