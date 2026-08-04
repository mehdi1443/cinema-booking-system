package dto

type UpdateUserProfileDTO struct {
	UserID    int64  `json:"user_id"  binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type UpdateUserRolesRequest struct {
	UserID  int64   `json:"user_id"  binding:"required"`
	RoleIDs []int64 `json:"role_ids" binding:"required"`
}

type HideFieldUserOption struct {
	WalletBalance bool
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

type AssignPermissionsToRoleRequest struct {
	PermissionIDs []*int64 `json:"permission_ids" binding:"required"`
}

type DischargePermissionsOfRoleRequest struct {
	PermissionIDs []*int64 `json:"permission_ids" binding:"required"`
}

//
//
//
//

// type CreateRoleDTO struct {
// 	Name  string
// 	Title string
// }

// type CreatePermissionDTO struct {
// 	Name        string
// 	Title       string
// 	Category    int
// 	Description string
// }

// type RegisterDTO struct {
// 	Phone     string
// 	Email     string
// 	Password  string
// 	FirstName string
// 	LastName  string
// }

// type LoginDTO struct {
// 	Identifier string // Can be phone or email
// 	Password   string
// }

// type AuthToken struct {
// 	AccessToken  string
// 	RefreshToken string
// }
