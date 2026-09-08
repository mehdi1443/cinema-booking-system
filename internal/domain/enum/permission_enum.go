package enum

import "strings"

type PermissionType string
type PermissionCategory uint

const (
	CategoryGeneral PermissionCategory = iota + 1
	CategoryUser
	CategoryMovie
	CategoryCinema
	CategoryInteraction
	CategoryRBAC
	CategoryLocation
	CategoryBooking
)

const (
	PermissionAll PermissionType = "general:all"
	AuthLogout    PermissionType = "auth:logout"

	UsersUpdatePassword     PermissionType = "users:update_password"
	UsersReadProfile        PermissionType = "users:read_profile"
	UsersReadProfileAny     PermissionType = "users:read_profile_any"
	UsersUpdateProfile      PermissionType = "users:update_profile"
	UsersChangeStatus       PermissionType = "users:change_status"
	UsersReadRoles          PermissionType = "users:read_roles"
	UsersReadRolesAny       PermissionType = "users:read_roles_any"
	UsersUpdateRoles        PermissionType = "users:update_roles"
	UsersReadPermissions    PermissionType = "users:read_permissions"
	UsersReadPermissionsAny PermissionType = "users:read_permissions_any"
	UsersDelete             PermissionType = "users:delete"
	UsersList               PermissionType = "users:list"
	WalletReadBalance       PermissionType = "wallet:read_balance"
	WalletReadBalanceAny    PermissionType = "wallet:read_balance_any"
	WalletChangeBalance     PermissionType = "wallet:change_balance"

	RolesRead              PermissionType = "roles:read"
	RolesCreate            PermissionType = "roles:create"
	RolesUpdate            PermissionType = "roles:update"
	RolesAppendPermissions PermissionType = "roles:append_permissions"
	RolesRevokePermissions PermissionType = "roles:revoke_permissions"
	PermissionsRead        PermissionType = "permissions:read"
	PermissionsReadByRole  PermissionType = "permissions:read_by_role"

	GenresCreate       PermissionType = "genres:create"
	GenresUpdate       PermissionType = "genres:update"
	GenresDelete       PermissionType = "genres:delete"
	PersonsCreate      PermissionType = "persons:create"
	PersonsUpdate      PermissionType = "persons:update"
	PersonsDelete      PermissionType = "persons:delete"
	PersonsManageImage PermissionType = "persons:manage_image"
	PersonsSyncSearch  PermissionType = "persons:sync_search"
	MoviesCreate       PermissionType = "movies:create"
	MoviesUpdate       PermissionType = "movies:update"
	MoviesDelete       PermissionType = "movies:delete"
	MoviesManageMedia  PermissionType = "movies:manage_media"
	MoviesSyncSearch   PermissionType = "movies:sync_search"

	ProvincesCreate  PermissionType = "provinces:create"
	ProvincesUpdate  PermissionType = "provinces:update"
	ProvincesDelete  PermissionType = "provinces:delete"
	CitiesCreate     PermissionType = "cities:create"
	CitiesUpdate     PermissionType = "cities:update"
	CitiesDelete     PermissionType = "cities:delete"
	FacilitiesCreate PermissionType = "facilities:create"
	FacilitiesUpdate PermissionType = "facilities:update"
	FacilitiesDelete PermissionType = "facilities:delete"

	CinemasCreate      PermissionType = "cinemas:create"
	CinemasUpdate      PermissionType = "cinemas:update"
	CinemasDelete      PermissionType = "cinemas:delete"
	CinemasManageImage PermissionType = "cinemas:manage_image"
	CinemasSyncSearch  PermissionType = "cinemas:sync_search"
	SearchSync         PermissionType = "search:sync"

	HallsCreate           PermissionType = "halls:create"
	HallsUpdate           PermissionType = "halls:update"
	HallsDelete           PermissionType = "halls:delete"
	HallsManageSeats      PermissionType = "halls:manage_seats"
	ShowtimesCreate       PermissionType = "showtimes:create"
	ShowtimesUpdate       PermissionType = "showtimes:update"
	ShowtimesDelete       PermissionType = "showtimes:delete"
	ShowtimesManagePrices PermissionType = "showtimes:manage_prices"
	ReservationsCreate    PermissionType = "reservations:create"
	ReservationsRead      PermissionType = "reservations:read"
	ReservationsCancel    PermissionType = "reservations:cancel"
	ReservationsReadAny   PermissionType = "reservations:read_any"
	ReservationsCancelAny PermissionType = "reservations:cancel_any"
	PaymentsCreate        PermissionType = "payments:create"
	PaymentsRead          PermissionType = "payments:read"
	PaymentsVerify        PermissionType = "payments:verify"
	PaymentsReadAny       PermissionType = "payments:read_any"
	PaymentsVerifyAny     PermissionType = "payments:verify_any"
	TicketsRead           PermissionType = "tickets:read"
	TicketsReadAny        PermissionType = "tickets:read_any"
	TicketsValidate       PermissionType = "tickets:validate"
	TicketsValidateAny    PermissionType = "tickets:validate_any"

	FavoritesAddMovie     PermissionType = "favorites:add_movie"
	FavoritesRemoveMovie  PermissionType = "favorites:remove_movie"
	FavoritesAddCinema    PermissionType = "favorites:add_cinema"
	FavoritesRemoveCinema PermissionType = "favorites:remove_cinema"
	FavoritesRead         PermissionType = "favorites:read"

	CommentsCreate       PermissionType = "comments:create"
	CommentsUpdate       PermissionType = "comments:update"
	CommentsDelete       PermissionType = "comments:delete"
	CommentsVote         PermissionType = "comments:vote"
	CommentsRemoveVote   PermissionType = "comments:remove_vote"
	CommentsManageStatus PermissionType = "comments:manage_status"
	CommentsAdminRead    PermissionType = "comments:admin_read"

	RatingsSubmit   PermissionType = "ratings:submit"
	RatingsReadSelf PermissionType = "ratings:read_self"
	RatingsDelete   PermissionType = "ratings:delete"
)

var allPermissionTypes = []PermissionType{
	PermissionAll, AuthLogout,
	UsersUpdatePassword, UsersReadProfile, UsersReadProfileAny, UsersUpdateProfile,
	UsersChangeStatus, UsersReadRoles, UsersReadRolesAny, UsersUpdateRoles,
	UsersReadPermissions, UsersReadPermissionsAny, UsersDelete, UsersList,
	WalletReadBalance, WalletReadBalanceAny, WalletChangeBalance,
	RolesRead, RolesCreate, RolesUpdate, RolesAppendPermissions, RolesRevokePermissions,
	PermissionsRead, PermissionsReadByRole,
	GenresCreate, GenresUpdate, GenresDelete,
	PersonsCreate, PersonsUpdate, PersonsDelete, PersonsManageImage, PersonsSyncSearch,
	MoviesCreate, MoviesUpdate, MoviesDelete, MoviesManageMedia, MoviesSyncSearch,
	ProvincesCreate, ProvincesUpdate, ProvincesDelete,
	CitiesCreate, CitiesUpdate, CitiesDelete,
	FacilitiesCreate, FacilitiesUpdate, FacilitiesDelete,
	CinemasCreate, CinemasUpdate, CinemasDelete, CinemasManageImage, CinemasSyncSearch, SearchSync,
	HallsCreate, HallsUpdate, HallsDelete, HallsManageSeats,
	ShowtimesCreate, ShowtimesUpdate, ShowtimesDelete, ShowtimesManagePrices,
	ReservationsCreate, ReservationsRead, ReservationsCancel, ReservationsReadAny, ReservationsCancelAny,
	PaymentsCreate, PaymentsRead, PaymentsVerify, PaymentsReadAny, PaymentsVerifyAny,
	TicketsRead, TicketsReadAny, TicketsValidate, TicketsValidateAny,
	FavoritesAddMovie, FavoritesRemoveMovie, FavoritesAddCinema, FavoritesRemoveCinema, FavoritesRead,
	CommentsCreate, CommentsUpdate, CommentsDelete, CommentsVote, CommentsRemoveVote, CommentsManageStatus, CommentsAdminRead,
	RatingsSubmit, RatingsReadSelf, RatingsDelete,
}

func (permission PermissionType) String() string { return string(permission) }
func (permission PermissionType) Description() string {
	return strings.ReplaceAll(string(permission), ":", " ")
}

func (permission PermissionType) Category() PermissionCategory {
	prefix, _, _ := strings.Cut(string(permission), ":")
	switch prefix {
	case "users", "wallet", "auth":
		return CategoryUser
	case "movies", "genres", "persons":
		return CategoryMovie
	case "cinemas", "facilities":
		return CategoryCinema
	case "halls", "showtimes", "reservations", "payments", "tickets":
		return CategoryBooking
	case "favorites", "comments", "ratings":
		return CategoryInteraction
	case "roles", "permissions":
		return CategoryRBAC
	case "provinces", "cities":
		return CategoryLocation
	default:
		return CategoryGeneral
	}
}

func (category PermissionCategory) String() string {
	switch category {
	case CategoryUser:
		return "user"
	case CategoryMovie:
		return "movie"
	case CategoryCinema:
		return "cinema"
	case CategoryInteraction:
		return "interaction"
	case CategoryRBAC:
		return "rbac"
	case CategoryLocation:
		return "location"
	case CategoryBooking:
		return "booking"
	default:
		return "general"
	}
}

func GetAllPermissionTypes() []PermissionType {
	result := make([]PermissionType, len(allPermissionTypes))
	copy(result, allPermissionTypes)
	return result
}

func GetUserPermissionTypes() []PermissionType {
	return []PermissionType{
		AuthLogout, UsersUpdatePassword, UsersReadProfile, UsersUpdateProfile,
		UsersReadRoles, UsersReadPermissions, WalletReadBalance,
		FavoritesAddMovie, FavoritesRemoveMovie, FavoritesAddCinema, FavoritesRemoveCinema, FavoritesRead,
		CommentsCreate, CommentsUpdate, CommentsDelete, CommentsVote, CommentsRemoveVote,
		RatingsSubmit, RatingsReadSelf, RatingsDelete,
		ReservationsCreate, ReservationsRead, ReservationsCancel,
		PaymentsCreate, PaymentsRead, PaymentsVerify,
		TicketsRead,
	}
}

func GetCinemaOwnerPermissionTypes() []PermissionType {
	permissions := GetUserPermissionTypes()
	return append(permissions,
		CinemasCreate, CinemasUpdate, CinemasDelete, CinemasManageImage,
		HallsCreate, HallsUpdate, HallsDelete, HallsManageSeats,
		ShowtimesCreate, ShowtimesUpdate, ShowtimesDelete, ShowtimesManagePrices,
		TicketsValidate,
	)
}

func GetSupportPermissionTypes() []PermissionType {
	permissions := GetUserPermissionTypes()
	return append(permissions, UsersList, UsersReadProfileAny, UsersReadRolesAny, UsersReadPermissionsAny, CommentsManageStatus, CommentsAdminRead, ReservationsReadAny, ReservationsCancelAny, PaymentsReadAny, TicketsReadAny, TicketsValidate, TicketsValidateAny)
}
