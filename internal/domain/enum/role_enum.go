package enum

type RoleName uint

const (
	RoleAdmin RoleName = iota + 1
	RoleCinemaOwner
	RoleSupport
	RoleCustomer
)

var roleNames = map[RoleName]string{
	RoleAdmin:       "admin",
	RoleCinemaOwner: "cinema_owner",
	RoleSupport:     "support",
	RoleCustomer:    "customer",
}

var roleTitles = map[RoleName]string{
	RoleAdmin:       "مدیر ارشد سیستم",
	RoleCinemaOwner: "مدیر سینما",
	RoleSupport:     "پشتیبان وب‌سایت",
	RoleCustomer:    "مشتری / کاربر عادی",
}

func (r RoleName) String() string {
	if name, ok := roleNames[r]; ok {
		return name
	}
	return "unknown"
}

func (r RoleName) Title() string {
	if title, ok := roleTitles[r]; ok {
		return title
	}
	return "ناشناخته"
}

func (r RoleName) DefaultPermissions() []PermissionType {
	switch r {
	case RoleAdmin:
		return []PermissionType{PermissionAll}

	case RoleCinemaOwner:
		return GetCinemaOwnerPermissionTypes()

	case RoleCustomer:
		return GetCustomerPermissionTypes()

	case RoleSupport:
		return []PermissionType{
			UserViewAll,
			TicketViewAll, TicketCancel,
			MovieViewAll, ShowtimeViewAll,
		}
	}
	return nil
}

func GetAllRoles() []RoleName {
	return []RoleName{
		RoleAdmin,
		RoleCinemaOwner,
		RoleSupport,
		RoleCustomer,
	}
}
