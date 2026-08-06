package service_contract

// import (
// 	"context"

// 	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
// 	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
// 	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
// )

// const (
// 	UserIDKeyMapClaims    string = "user_id"
// 	TokenTypeKeyMapClaims string = "token_type"
// 	IssuedAtKeyMapClaims  string = "iat"
// 	ExpireAtKeyMapClaims  string = "exp"
// )

// type (
// 	JWTServiceInterface interface {
// 		GenerateAccessToken(req dto.GenerateTokenRequest) (string, error)
// 		GenerateRefreshToken(req dto.GenerateTokenRequest) (string, error)
// 		ValidateToken(signedToken string) (*domainjwt.TokenPayload, error)
// 	}

// 	UserServiceInterface interface {
// 		GetUserByID(ctx context.Context, id int64, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error)
// 		GetUserByPhone(ctx context.Context, phone string, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error)
// 		GetUserByEmail(ctx context.Context, email string, hideOpt *dto.HideFieldUserOption) (*dto.UserResponse, error)

// 		GetUsersByStatus(ctx context.Context, status enum.UserStatus, hideOpt *dto.HideFieldUserOption) ([]*dto.UserResponse, error)

// 		UpdateProfile(ctx context.Context, dto dto.UpdateUserProfileDTO) (*dto.UserResponse, error)
// 		ChangeStatus(ctx context.Context, userID int64, status enum.UserStatus) (*dto.UserResponse, error)
// 		UpdateUserRoles(ctx context.Context, request dto.UpdateUserRolesRequest) (*dto.UserResponse, error)

// 		DeleteUserByID(ctx context.Context, id int64) error
// 	}

// 	WalletServiceInterface interface {
// 		GetBalance(ctx context.Context, userID int64) (int64, error)
// 		IncreaseBallance(ctx context.Context, userID int64, amount int64) (int64, error)
// 		DecreaseBalance(ctx context.Context, userID int64, amount int64) (int64, error)
// 	}

// 	RBACServiceInterface interface {
// 		// Appendments & Discharging
// 		AppendPermissionsToRole(ctx context.Context, roleID int64, permissionID []*int64) error
// 		DeletePermissionsOfRole(ctx context.Context, roleID int64, permissionID []*int64) error

// 		// Authorization Enforcement
// 		UserHasPermission(ctx context.Context, userID int64, permissionName string) (bool, error)
// 		UserHasRole(ctx context.Context, userID int64, roleName string) (bool, error)

// 		GetUserPermissions(ctx context.Context, userID int64) ([]*dto.PermissionResponse, error)
// 		GetUserRoles(ctx context.Context, userID int64) ([]*dto.RoleResponse, error)

// 		GetUsersByPermissionName(ctx context.Context, permissionNames []string) ([]*dto.UserResponse, error)
// 		GetUsersByRoleNames(ctx context.Context, roleNames []string) ([]*dto.UserResponse, error)

// 		GetPermissionsByRoleName(ctx context.Context, role string) ([]*dto.PermissionResponse, error)
// 		GetRolesByPermissionName(ctx context.Context, permissionName []string) ([]*dto.RoleResponse, error)

// 		GetPermissionsByID(ctx context.Context, pid int64) (*dto.PermissionResponse, error)
// 		GetRoleByID(ctx context.Context, rid int64) (*dto.RoleResponse, error)
// 	}

// 	// AuthenticationServiceInterface interface {
// 	// 	ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) error
// 	// 	VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) (userdto.UserInfoResponse, error)
// 	// 	Register(registerInfo userdto.BasicRegisterRequest) error
// 	// 	VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) error
// 	// 	Login(loginInfo userdto.LoginRequest) (userdto.UserInfoResponse, error)
// 	// 	CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) error

// 	// 	VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest) error
// 	// 	ResetPassword(resetPassInfo userdto.ResetPasswordRequest) error

// 	// 	Register(ctx context.Context, dto RegisterDTO) (*entity.User, *AuthToken, error)
// 	// 	Login(ctx context.Context, dto LoginDTO) (*entity.User, *AuthToken, error)
// 	// 	RefreshToken(ctx context.Context, refreshToken string) (*AuthToken, error)
// 	// 	RefreshToken(refreshToken string) (userdto.UserInfoResponse, error)
// 	// }
// )
