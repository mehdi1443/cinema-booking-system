package enum

type PermissionType uint
type PermissionCategory uint

const (
	PermissionAll PermissionType = iota + 1
	PermissionGeneral

	// Admin / User Management
	UserViewAll
	UserBanUnban
	UserChangeRole
	UserViewRoles
	UserManageRolePermissions

	// Movie Management
	MovieViewAll
	MovieCreate
	MovieEdit
	MovieDelete

	// Cinema & Hall Management
	CinemaViewAll
	CinemaCreate
	CinemaEdit
	CinemaDelete
	HallCreate
	HallEdit
	HallDelete

	// Showtime Management
	ShowtimeViewAll
	ShowtimeCreate
	ShowtimeEdit
	ShowtimeDelete

	// Ticket & Reservation Management
	TicketBuy
	TicketCancel
	TicketViewOwn
	TicketViewAll

	// Wallet & Financial Management
	WalletCharge
	WalletViewBalance
	FinancialReportView
)

const (
	CategoryGeneral PermissionCategory = iota + 1
	CategoryUser
	CategoryMovie
	CategoryCinema
	CategoryShowtime
	CategoryTicket
	CategoryWallet
)

var permissionNames = map[PermissionType]string{
	PermissionAll:     "general.all",
	PermissionGeneral: "general.general",

	// User Management
	UserViewAll:               "user.viewAll",
	UserBanUnban:              "user.banUnban",
	UserChangeRole:            "user.changeRole",
	UserViewRoles:             "user.viewRoles",
	UserManageRolePermissions: "user.manageRolePermissions",

	// Movie Management
	MovieViewAll: "movie.viewAll",
	MovieCreate:  "movie.create",
	MovieEdit:    "movie.edit",
	MovieDelete:  "movie.delete",

	// Cinema & Hall Management
	CinemaViewAll: "cinema.viewAll",
	CinemaCreate:  "cinema.create",
	CinemaEdit:    "cinema.edit",
	CinemaDelete:  "cinema.delete",
	HallCreate:    "hall.create",
	HallEdit:      "hall.edit",
	HallDelete:    "hall.delete",

	// Showtime Management
	ShowtimeViewAll: "showtime.viewAll",
	ShowtimeCreate:  "showtime.create",
	ShowtimeEdit:    "showtime.edit",
	ShowtimeDelete:  "showtime.delete",

	// Ticket Management
	TicketBuy:     "ticket.buy",
	TicketCancel:  "ticket.cancel",
	TicketViewOwn: "ticket.viewOwn",
	TicketViewAll: "ticket.viewAll",

	// Wallet Management
	WalletCharge:        "wallet.charge",
	WalletViewBalance:   "wallet.viewBalance",
	FinancialReportView: "financial.reportView",
}

var permissionDescriptions = map[PermissionType]string{
	PermissionAll:     "دسترسی کامل به کل سیستم سینما",
	PermissionGeneral: "دسترسی عمومی",

	// User Management
	UserViewAll:               "مشاهده لیست کاربران",
	UserBanUnban:              "مسدود/رفع مسدودیت کاربران",
	UserChangeRole:            "تغییر نقش کاربر",
	UserViewRoles:             "مشاهده لیست نقش‌ها و دسترسی‌ها",
	UserManageRolePermissions: "تغییر دسترسی‌های نقش",

	// Movie Management
	MovieViewAll: "مشاهده لیست تمام فیلم‌ها",
	MovieCreate:  "ایجاد فیلم جدید",
	MovieEdit:    "ویرایش اطلاعات فیلم",
	MovieDelete:  "حذف فیلم",

	// Cinema & Hall Management
	CinemaViewAll: "مشاهده لیست سینماها",
	CinemaCreate:  "تعریف سینمای جدید",
	CinemaEdit:    "ویرایش اطلاعات سینما",
	CinemaDelete:  "حذف سینما",
	HallCreate:    "افزودن سالن به سینما",
	HallEdit:      "ویرایش سالن و چیدمان صندلی‌ها",
	HallDelete:    "حذف سالن سینما",

	// Showtime Management
	ShowtimeViewAll: "مشاهده جدول سانس‌ها",
	ShowtimeCreate:  "تعریف سانس جدید برای اکران",
	ShowtimeEdit:    "ویرایش زمان و قیمت سانس",
	ShowtimeDelete:  "لغو/حذف سانس",

	// Ticket Management
	TicketBuy:     "خرید و رزرو بلیت",
	TicketCancel:  "استرداد و لغو بلیت",
	TicketViewOwn: "مشاهده بلیت‌های شخصی",
	TicketViewAll: "مشاهده تمام بلیت‌های صادر شده",

	// Wallet Management
	WalletCharge:        "شارژ کیف پول",
	WalletViewBalance:   "مشاهده موجودی کیف پول",
	FinancialReportView: "مشاهده گزارش‌های مالی فروش",
}

var permissionCategories = map[PermissionType]PermissionCategory{
	PermissionAll:     CategoryGeneral,
	PermissionGeneral: CategoryGeneral,

	// User Management
	UserViewAll:               CategoryUser,
	UserBanUnban:              CategoryUser,
	UserChangeRole:            CategoryUser,
	UserViewRoles:             CategoryUser,
	UserManageRolePermissions: CategoryUser,

	// Movie Management
	MovieViewAll: CategoryMovie,
	MovieCreate:  CategoryMovie,
	MovieEdit:    CategoryMovie,
	MovieDelete:  CategoryMovie,

	// Cinema & Hall Management
	CinemaViewAll: CategoryCinema,
	CinemaCreate:  CategoryCinema,
	CinemaEdit:    CategoryCinema,
	CinemaDelete:  CategoryCinema,
	HallCreate:    CategoryCinema,
	HallEdit:      CategoryCinema,
	HallDelete:    CategoryCinema,

	// Showtime Management
	ShowtimeViewAll: CategoryShowtime,
	ShowtimeCreate:  CategoryShowtime,
	ShowtimeEdit:    CategoryShowtime,
	ShowtimeDelete:  CategoryShowtime,

	// Ticket Management
	TicketBuy:     CategoryTicket,
	TicketCancel:  CategoryTicket,
	TicketViewOwn: CategoryTicket,
	TicketViewAll: CategoryTicket,

	// Wallet Management
	WalletCharge:        CategoryWallet,
	WalletViewBalance:   CategoryWallet,
	FinancialReportView: CategoryWallet,
}

func (perm PermissionType) String() string {
	if name, ok := permissionNames[perm]; ok {
		return name
	}
	return ""
}

func (perm PermissionType) Description() string {
	if desc, ok := permissionDescriptions[perm]; ok {
		return desc
	}
	return ""
}

func (perm PermissionType) Category() PermissionCategory {
	if category, ok := permissionCategories[perm]; ok {
		return category
	}
	return CategoryGeneral
}

func (category PermissionCategory) String() string {
	switch category {
	case CategoryGeneral:
		return "عمومی"
	case CategoryUser:
		return "مدیریت کاربران"
	case CategoryMovie:
		return "مدیریت فیلم‌ها"
	case CategoryCinema:
		return "مدیریت سینما و سالن‌ها"
	case CategoryShowtime:
		return "مدیریت سانس‌ها"
	case CategoryTicket:
		return "مدیریت بلیت‌ها"
	case CategoryWallet:
		return "کیف پول و امور مالی"
	}
	return "unknown"
}

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		PermissionAll, PermissionGeneral,
		UserViewAll, UserBanUnban, UserChangeRole, UserViewRoles, UserManageRolePermissions,
		MovieViewAll, MovieCreate, MovieEdit, MovieDelete,
		CinemaViewAll, CinemaCreate, CinemaEdit, CinemaDelete, HallCreate, HallEdit, HallDelete,
		ShowtimeViewAll, ShowtimeCreate, ShowtimeEdit, ShowtimeDelete,
		TicketBuy, TicketCancel, TicketViewOwn, TicketViewAll,
		WalletCharge, WalletViewBalance, FinancialReportView,
	}
}

func GetCinemaOwnerPermissionTypes() []PermissionType {
	return []PermissionType{
		CinemaViewAll, CinemaEdit, HallCreate, HallEdit, HallDelete,
		ShowtimeViewAll, ShowtimeCreate, ShowtimeEdit, ShowtimeDelete,
		TicketViewAll, FinancialReportView,
	}
}

func GetCustomerPermissionTypes() []PermissionType {
	return []PermissionType{
		MovieViewAll, ShowtimeViewAll,
		TicketBuy, TicketCancel, TicketViewOwn,
		WalletCharge, WalletViewBalance,
	}
}
