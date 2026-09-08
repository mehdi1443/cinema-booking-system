package enum

import "testing"

func containsPermission(permissions []PermissionType, expected PermissionType) bool {
	for _, permission := range permissions {
		if permission == expected {
			return true
		}
	}
	return false
}

func TestBuiltInRoleNamesAndPermissions(t *testing.T) {
	t.Parallel()
	if RoleAdmin.String() != "admin" || RoleSupport.String() != "support" || RoleUser.String() != "user" || RoleCinemaOwner.String() != "cinema_owner" {
		t.Fatal("built-in role names changed")
	}
	if !containsPermission(RoleAdmin.DefaultPermissions(), PermissionAll) {
		t.Fatal("admin must have the wildcard permission")
	}
	if containsPermission(RoleUser.DefaultPermissions(), WalletChangeBalance) {
		t.Fatal("ordinary users must not directly change wallet balances")
	}
	if !containsPermission(RoleCinemaOwner.DefaultPermissions(), CinemasUpdate) {
		t.Fatal("cinema owners must be able to update owned cinemas")
	}
	for _, permission := range []PermissionType{HallsManageSeats, ShowtimesCreate, ShowtimesManagePrices, ReservationsCreate} {
		if !containsPermission(RoleCinemaOwner.DefaultPermissions(), permission) {
			t.Fatalf("cinema owner missing permission %s", permission)
		}
	}
	if !containsPermission(RoleSupport.DefaultPermissions(), CommentsManageStatus) {
		t.Fatal("support must be able to moderate comments")
	}
	if !containsPermission(RoleSupport.DefaultPermissions(), ReservationsReadAny) {
		t.Fatal("support must be able to inspect reservations")
	}
}
