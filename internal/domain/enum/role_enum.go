package enum

type RoleName uint

const (
	RoleAdmin RoleName = iota + 1
	RoleCinemaOwner
	RoleSupport
	RoleUser
)

var roleNames = map[RoleName]string{
	RoleAdmin: "admin", RoleCinemaOwner: "cinema_owner", RoleSupport: "support", RoleUser: "user",
}

var roleTitles = map[RoleName]string{
	RoleAdmin: "مدیر سیستم", RoleCinemaOwner: "صاحب سینما", RoleSupport: "پشتیبان", RoleUser: "کاربر عادی",
}

func (role RoleName) String() string { return roleNames[role] }
func (role RoleName) Title() string  { return roleTitles[role] }

func (role RoleName) DefaultPermissions() []PermissionType {
	switch role {
	case RoleAdmin:
		return []PermissionType{PermissionAll}
	case RoleCinemaOwner:
		return GetCinemaOwnerPermissionTypes()
	case RoleSupport:
		return GetSupportPermissionTypes()
	case RoleUser:
		return GetUserPermissionTypes()
	default:
		return nil
	}
}

func GetAllRoles() []RoleName {
	return []RoleName{RoleAdmin, RoleCinemaOwner, RoleSupport, RoleUser}
}
