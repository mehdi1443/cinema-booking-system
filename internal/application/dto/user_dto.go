package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type HideFieldUserOption struct {
	WalletBalance bool
}

type UserResponse struct {
	ID            int64           `json:"id"`
	Phone         string          `json:"phone,omitempty"`
	Email         string          `json:"email,omitempty"`
	FirstName     string          `json:"first_name,omitempty"`
	LastName      string          `json:"last_name,omitempty"`
	Status        enum.UserStatus `json:"status,omitempty"`
	WalletBalance *int64          `json:"wallet_balance,omitempty"`
	CreatedAt     *time.Time      `json:"created_at,omitempty"`
	UpdatedAt     *time.Time      `json:"updated_at,omitempty"`
	DeletedAt     *time.Time      `json:"deleted_at,omitempty"`
	Roles         []*RoleResponse `json:"roles,omitempty"`

	IsPhoneVerified bool `json:"is_verified_phone"`
	IsEmailVerified bool `json:"is_verified_email"`
}

type ChangeUserStatusRequest struct {
	Status enum.UserStatus `json:"status" binding:"required"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type GetUsersQueryDTO struct {
	Phone       string   `form:"phone"`
	Email       string   `form:"email"`
	FirstName   string   `form:"first_name"`
	LastName    string   `form:"last_name"`
	Status      string   `form:"status"`
	Permissions []string `form:"permission"`
	Roles       []string `form:"role"`
	HideFieldUserOption
}

// IsEmpty checks if at least one filter criterion is provided
func (q *GetUsersQueryDTO) IsEmpty() bool {
	return q.Phone == "" &&
		q.Email == "" &&
		q.FirstName == "" &&
		q.LastName == "" &&
		q.Status == "" &&
		len(q.Permissions) == 0 &&
		len(q.Roles) == 0
}

func ToUserResponse(user *entity.User, hideOpt *HideFieldUserOption) *UserResponse {
	if user == nil {
		return nil
	}

	res := &UserResponse{
		ID:              user.ID,
		Phone:           user.Phone,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Status:          user.Status,
		WalletBalance:   &user.WalletBalance,
		CreatedAt:       helper.TimeToPtr(user.CreatedAt),
		UpdatedAt:       helper.TimeToPtr(user.UpdatedAt),
		IsPhoneVerified: user.IsPhoneVerified,
		IsEmailVerified: user.IsEmailVerified,
	}

	if hideOpt != nil {
		res.WalletBalance = helper.HideValue[int64](user.WalletBalance, hideOpt.WalletBalance)
	}

	// if !user.DeletedAt.IsZero() {
	// 	res.DeletedAt = helper.TimeToPtr(user.DeletedAt)
	// }

	// Only populate roles if relationship was preloaded
	if len(user.Roles) > 0 {
		res.Roles = ToRoleResponseSlice(user.Roles)
	}

	return res
}

func ToUserResponseSlice(users []*entity.User, hideOpt *HideFieldUserOption) []*UserResponse {
	if len(users) == 0 {
		return []*UserResponse{}
	}

	responses := make([]*UserResponse, 0, len(users))
	for ind := range users {
		if resp := ToUserResponse(users[ind], hideOpt); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}
