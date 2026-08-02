package dto

type UpdateUserProfileDTO struct {
	UserID    int64
	FirstName string
	LastName  string
	Email     string
}

type UpdateUserRolesRequest struct {
	UserID  int64
	RoleIDs []int64
}

type HideFieldUserOption struct {
	WalletBalance bool
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
