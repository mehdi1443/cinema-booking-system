package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type endpointPolicy struct {
	permission      string
	selfPermission  string
	anyPermission   string
	userIDParam     string
	ownedCinema     bool
	ownedHall       bool
	ownedSeat       bool
	ownedShowtime   bool
	resourceIDParam string
}

func permissionPolicy(permission enum.PermissionType) endpointPolicy {
	return endpointPolicy{permission: permission.String()}
}

func selfPolicy(selfPermission, anyPermission enum.PermissionType) endpointPolicy {
	return endpointPolicy{selfPermission: selfPermission.String(), anyPermission: anyPermission.String(), userIDParam: "id"}
}

func ownedCinemaPolicy(permission enum.PermissionType) endpointPolicy {
	return endpointPolicy{permission: permission.String(), ownedCinema: true, resourceIDParam: "id"}
}

func ownedCinemaParamPolicy(permission enum.PermissionType, param string) endpointPolicy {
	return endpointPolicy{permission: permission.String(), ownedCinema: true, resourceIDParam: param}
}

func ownedHallPolicy(permission enum.PermissionType) endpointPolicy {
	return endpointPolicy{permission: permission.String(), ownedHall: true, resourceIDParam: "id"}
}

func ownedSeatPolicy(permission enum.PermissionType) endpointPolicy {
	return endpointPolicy{permission: permission.String(), ownedSeat: true, resourceIDParam: "id"}
}

func ownedShowtimePolicy(permission enum.PermissionType) endpointPolicy {
	return endpointPolicy{permission: permission.String(), ownedShowtime: true, resourceIDParam: "id"}
}

var publicEndpointPolicies = map[string]struct{}{
	"POST /v1/auth/otp/send": {}, "POST /v1/auth/otp/verify": {}, "POST /v1/auth/refresh": {},
	"POST /v1/captcha/send": {}, "POST /v1/captcha/verify": {},
	"GET /v1/genres": {}, "GET /v1/genres/:id": {},
	"GET /v1/persons": {}, "GET /v1/persons/:id": {},
	"GET /v1/movies": {}, "GET /v1/movies/:id": {},
	"GET /v1/provinces": {}, "GET /v1/provinces/:id": {},
	"GET /v1/cities": {}, "GET /v1/cities/:id": {},
	"GET /v1/facilities": {}, "GET /v1/facilities/:id": {},
	"GET /v1/cinemas": {}, "GET /v1/cinemas/:id": {},
	"GET /v1/comments/:id": {}, "GET /v1/comments/:id/replies": {},
	"GET /v1/movies/:id/comments": {}, "GET /v1/cinemas/:id/comments": {},
	"GET /v1/movies/:id/ratings": {}, "GET /v1/cinemas/:id/ratings": {},
	"GET /v1/search/": {}, "GET /v1/search/suggest": {},
	"GET /v1/cinemas/:id/halls": {}, "GET /v1/halls/:id": {}, "GET /v1/halls/:id/seats": {},
	"GET /v1/showtimes": {}, "GET /v1/showtimes/:id": {}, "GET /v1/showtimes/:id/seats": {},
}

var protectedEndpointPolicies = map[string]endpointPolicy{
	"POST /v1/auth/logout": permissionPolicy(enum.AuthLogout),

	"PATCH /v1/users/:id/password":      selfPolicy(enum.UsersUpdatePassword, enum.PermissionAll),
	"GET /v1/users/:id/profile":         selfPolicy(enum.UsersReadProfile, enum.UsersReadProfileAny),
	"PUT /v1/users/:id/profile":         selfPolicy(enum.UsersUpdateProfile, enum.PermissionAll),
	"PATCH /v1/users/:id/status":        permissionPolicy(enum.UsersChangeStatus),
	"GET /v1/users/:id/roles":           selfPolicy(enum.UsersReadRoles, enum.UsersReadRolesAny),
	"PUT /v1/users/:id/roles":           permissionPolicy(enum.UsersUpdateRoles),
	"GET /v1/users/:id/permissions":     selfPolicy(enum.UsersReadPermissions, enum.UsersReadPermissionsAny),
	"DELETE /v1/users/:id":              permissionPolicy(enum.UsersDelete),
	"GET /v1/users":                     permissionPolicy(enum.UsersList),
	"GET /v1/users/:id/wallet/balance":  selfPolicy(enum.WalletReadBalance, enum.WalletReadBalanceAny),
	"POST /v1/users/:id/wallet/balance": permissionPolicy(enum.WalletChangeBalance),

	"GET /v1/roles/:id":                permissionPolicy(enum.RolesRead),
	"POST /v1/roles/":                  permissionPolicy(enum.RolesCreate),
	"PUT /v1/roles/:id":                permissionPolicy(enum.RolesUpdate),
	"POST /v1/roles/:id/permissions":   permissionPolicy(enum.RolesAppendPermissions),
	"DELETE /v1/roles/:id/permissions": permissionPolicy(enum.RolesRevokePermissions),
	"GET /v1/permissions/:id":          permissionPolicy(enum.PermissionsRead),
	"GET /v1/permissions/role/:name":   permissionPolicy(enum.PermissionsReadByRole),

	"POST /v1/genres": permissionPolicy(enum.GenresCreate), "PUT /v1/genres/:id": permissionPolicy(enum.GenresUpdate), "DELETE /v1/genres/:id": permissionPolicy(enum.GenresDelete),
	"POST /v1/persons": permissionPolicy(enum.PersonsCreate), "PUT /v1/persons/:id": permissionPolicy(enum.PersonsUpdate), "DELETE /v1/persons/:id": permissionPolicy(enum.PersonsDelete),
	"POST /v1/persons/:id/image": permissionPolicy(enum.PersonsManageImage), "DELETE /v1/persons/:id/image": permissionPolicy(enum.PersonsManageImage), "POST /v1/persons/sync-search": permissionPolicy(enum.PersonsSyncSearch),
	"POST /v1/movies": permissionPolicy(enum.MoviesCreate), "PUT /v1/movies/:id": permissionPolicy(enum.MoviesUpdate), "DELETE /v1/movies/:id": permissionPolicy(enum.MoviesDelete),
	"POST /v1/movies/:id/poster": permissionPolicy(enum.MoviesManageMedia), "DELETE /v1/movies/:id/poster": permissionPolicy(enum.MoviesManageMedia),
	"POST /v1/movies/:id/trailer": permissionPolicy(enum.MoviesManageMedia), "DELETE /v1/movies/:id/trailer": permissionPolicy(enum.MoviesManageMedia), "POST /v1/movies/sync-search": permissionPolicy(enum.MoviesSyncSearch),
	"POST /v1/provinces": permissionPolicy(enum.ProvincesCreate), "PUT /v1/provinces/:id": permissionPolicy(enum.ProvincesUpdate), "DELETE /v1/provinces/:id": permissionPolicy(enum.ProvincesDelete),
	"POST /v1/cities": permissionPolicy(enum.CitiesCreate), "PUT /v1/cities/:id": permissionPolicy(enum.CitiesUpdate), "DELETE /v1/cities/:id": permissionPolicy(enum.CitiesDelete),
	"POST /v1/facilities": permissionPolicy(enum.FacilitiesCreate), "PUT /v1/facilities/:id": permissionPolicy(enum.FacilitiesUpdate), "DELETE /v1/facilities/:id": permissionPolicy(enum.FacilitiesDelete),
	"POST /v1/cinemas": permissionPolicy(enum.CinemasCreate), "PUT /v1/cinemas/:id": ownedCinemaPolicy(enum.CinemasUpdate), "DELETE /v1/cinemas/:id": ownedCinemaPolicy(enum.CinemasDelete),
	"POST /v1/cinemas/:id/image": ownedCinemaPolicy(enum.CinemasManageImage), "DELETE /v1/cinemas/:id/image": ownedCinemaPolicy(enum.CinemasManageImage), "POST /v1/cinemas/sync-search": permissionPolicy(enum.CinemasSyncSearch),
	"POST /v1/search/sync": permissionPolicy(enum.SearchSync),

	"POST /v1/favorites/movies": permissionPolicy(enum.FavoritesAddMovie), "DELETE /v1/favorites/movies/:movie_id": permissionPolicy(enum.FavoritesRemoveMovie),
	"POST /v1/favorites/cinemas": permissionPolicy(enum.FavoritesAddCinema), "DELETE /v1/favorites/cinemas/:cinema_id": permissionPolicy(enum.FavoritesRemoveCinema), "GET /v1/favorites": permissionPolicy(enum.FavoritesRead),
	"POST /v1/comments": permissionPolicy(enum.CommentsCreate), "PUT /v1/comments/:id": permissionPolicy(enum.CommentsUpdate), "DELETE /v1/comments/:id": permissionPolicy(enum.CommentsDelete),
	"POST /v1/comments/:id/vote": permissionPolicy(enum.CommentsVote), "DELETE /v1/comments/:id/vote": permissionPolicy(enum.CommentsRemoveVote), "PATCH /v1/comments/:id/status": permissionPolicy(enum.CommentsManageStatus),
	"GET /v1/admin/comments": permissionPolicy(enum.CommentsAdminRead), "GET /v1/admin/movies/:id/comments": permissionPolicy(enum.CommentsAdminRead),
	"GET /v1/admin/cinemas/:id/comments": permissionPolicy(enum.CommentsAdminRead), "GET /v1/admin/comments/:id/replies": permissionPolicy(enum.CommentsAdminRead),
	"POST /v1/ratings": permissionPolicy(enum.RatingsSubmit), "GET /v1/movies/:id/user-rating": permissionPolicy(enum.RatingsReadSelf),
	"GET /v1/cinemas/:id/user-rating": permissionPolicy(enum.RatingsReadSelf), "DELETE /v1/ratings/:id": permissionPolicy(enum.RatingsDelete),

	"POST /v1/cinemas/:id/halls": ownedCinemaParamPolicy(enum.HallsCreate, "id"),
	"PUT /v1/halls/:id":          ownedHallPolicy(enum.HallsUpdate), "DELETE /v1/halls/:id": ownedHallPolicy(enum.HallsDelete),
	"POST /v1/halls/:id/seats": ownedHallPolicy(enum.HallsManageSeats), "POST /v1/halls/:id/showtimes": ownedHallPolicy(enum.ShowtimesCreate),
	"PUT /v1/seats/:id": ownedSeatPolicy(enum.HallsManageSeats), "DELETE /v1/seats/:id": ownedSeatPolicy(enum.HallsManageSeats),
	"PUT /v1/showtimes/:id": ownedShowtimePolicy(enum.ShowtimesUpdate), "DELETE /v1/showtimes/:id": ownedShowtimePolicy(enum.ShowtimesDelete),
	"PUT /v1/showtimes/:id/prices": ownedShowtimePolicy(enum.ShowtimesManagePrices),
	"POST /v1/reservations":        permissionPolicy(enum.ReservationsCreate), "GET /v1/reservations": permissionPolicy(enum.ReservationsRead),
	"GET /v1/reservations/:id": permissionPolicy(enum.ReservationsRead), "DELETE /v1/reservations/:id": permissionPolicy(enum.ReservationsCancel),
	"POST /v1/reservations/:id/payments": permissionPolicy(enum.PaymentsCreate), "GET /v1/reservations/:id/payments": permissionPolicy(enum.PaymentsRead),
	"GET /v1/payments/:id": permissionPolicy(enum.PaymentsRead), "POST /v1/payments/:id/verify": permissionPolicy(enum.PaymentsVerify),
	"GET /v1/tickets": permissionPolicy(enum.TicketsRead), "GET /v1/tickets/:id": permissionPolicy(enum.TicketsRead),
	"GET /v1/tickets/:id/qr": permissionPolicy(enum.TicketsRead), "POST /v1/tickets/validate": permissionPolicy(enum.TicketsValidate),
}

// ValidateRoutePolicies makes authorization configuration fail fast when a
// registered v1 endpoint has no explicit public or protected policy. It also
// detects stale policy entries and accidental public/protected overlaps.
func ValidateRoutePolicies(routes gin.RoutesInfo) error {
	registered := make(map[string]struct{}, len(routes))
	problems := make([]string, 0)

	for _, route := range routes {
		if !strings.HasPrefix(route.Path, "/v1/") {
			continue
		}
		key := route.Method + " " + route.Path
		registered[key] = struct{}{}
		_, public := publicEndpointPolicies[key]
		_, protected := protectedEndpointPolicies[key]
		switch {
		case public && protected:
			problems = append(problems, "policy is both public and protected: "+key)
		case !public && !protected:
			problems = append(problems, "missing policy: "+key)
		}
	}

	for key := range publicEndpointPolicies {
		if _, ok := registered[key]; !ok {
			problems = append(problems, "public policy has no registered route: "+key)
		}
	}
	for key := range protectedEndpointPolicies {
		if _, ok := registered[key]; !ok {
			problems = append(problems, "protected policy has no registered route: "+key)
		}
	}

	if len(problems) == 0 {
		return nil
	}
	sort.Strings(problems)
	return fmt.Errorf("route authorization policy mismatch: %s", strings.Join(problems, "; "))
}

func RouteAuthorization(
	tokenService service_contract.JWTServiceInterface,
	rbacRepo DB.RBACInterface,
	permissionCache cache.UserPermissionCacheInterface,
	contextTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	cinemaRepo DB.CinemaInterface,
	bookingRepo DB.BookingRepositoryInterface,
	cacheConfig *bootstrap.UserPermissionCache,
	log domainlogger.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			c.Next()
			return
		}
		key := c.Request.Method + " " + path
		if _, public := publicEndpointPolicies[key]; public {
			c.Next()
			return
		}

		policy, known := protectedEndpointPolicies[key]
		if !known {
			policy = permissionPolicy(enum.PermissionAll)
		}
		if !authenticateRoute(c, tokenService, log) {
			return
		}
		userID, permissions, ok := loadUserPermissions(c, rbacRepo, permissionCache, contextTokenService, userRepo, cacheConfig, log)
		if !ok {
			return
		}
		requiredPermission := policy.permission
		if policy.userIDParam != "" {
			targetID, err := strconv.ParseInt(c.Param(policy.userIDParam), 10, 64)
			if err != nil || targetID <= 0 {
				response.ErrorWithCode(c, http.StatusBadRequest, string(exception.CodeInvalidInput), "Invalid user id")
				c.Abort()
				return
			}
			requiredPermission = policy.anyPermission
			if targetID == userID {
				requiredPermission = policy.selfPermission
			}
		}
		if !hasPermission(permissions, requiredPermission) {
			response.ErrorWithCode(c, http.StatusForbidden, string(exception.CodeForbidden), "You do not have permission to access this resource")
			c.Abort()
			return
		}
		if (policy.ownedCinema || policy.ownedHall || policy.ownedSeat || policy.ownedShowtime) && !hasPermission(permissions, enum.PermissionAll.String()) {
			resourceID, err := strconv.ParseInt(c.Param(policy.resourceIDParam), 10, 64)
			if err != nil || resourceID <= 0 {
				response.ErrorWithCode(c, http.StatusBadRequest, string(exception.CodeInvalidInput), "Invalid resource id")
				c.Abort()
				return
			}
			var owned bool
			switch {
			case policy.ownedCinema:
				owned, err = cinemaRepo.IsOwnedBy(c.Request.Context(), resourceID, userID)
			case policy.ownedHall:
				owned, err = bookingRepo.IsHallOwnedBy(c.Request.Context(), resourceID, userID)
			case policy.ownedSeat:
				owned, err = bookingRepo.IsSeatOwnedBy(c.Request.Context(), resourceID, userID)
			case policy.ownedShowtime:
				owned, err = bookingRepo.IsShowtimeOwnedBy(c.Request.Context(), resourceID, userID)
			}
			if err != nil {
				log.Error(c.Request.Context(), "failed to verify resource ownership",
					domainlogger.Int64("user_id", userID),
					domainlogger.Int64("resource_id", resourceID),
					domainlogger.Error(err),
				)
				response.ErrorWithCode(c, http.StatusInternalServerError, string(exception.CodeInternal), "Failed to verify resource ownership")
				c.Abort()
				return
			}
			if !owned {
				response.ErrorWithCode(c, http.StatusForbidden, string(exception.CodeForbidden), "You can only manage resources belonging to your cinema")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func authenticateRoute(c *gin.Context, tokenService service_contract.JWTServiceInterface, log domainlogger.Logger) bool {
	authorization := c.GetHeader("Authorization")
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		response.ErrorWithCode(c, http.StatusUnauthorized, string(exception.CodeUnauthorized), "Authorization header is required")
		c.Abort()
		return false
	}
	payload, err := tokenService.ValidateToken(c.Request.Context(), parts[1], string(domainjwt.AccessTokenTag))
	if err != nil {
		if errors.Is(err, exception.ErrTokenExpired) {
			response.ErrorWithCode(c, http.StatusUnauthorized, string(exception.CodeTokenExpired), "access token has expired")
		} else {
			log.Warn(c.Request.Context(), "failed token validation", domainlogger.Error(err))
			response.ErrorWithCode(c, http.StatusUnauthorized, string(exception.CodeInvalidToken), "token is malformed or invalid")
		}
		c.Abort()
		return false
	}
	if payload.TokenType != domainjwt.AccessTokenTag {
		response.ErrorWithCode(c, http.StatusForbidden, string(exception.CodeForbidden), "Invalid token type")
		c.Abort()
		return false
	}
	blacklisted, err := tokenService.IsBlacklisted(c.Request.Context(), parts[1])
	if err != nil {
		log.Error(c.Request.Context(), "failed to check token revocation status", domainlogger.Int64("user_id", payload.UserID), domainlogger.Error(err))
		response.ErrorWithCode(c, http.StatusInternalServerError, string(exception.CodeInternal), "Failed to verify token status")
		c.Abort()
		return false
	}
	if blacklisted != nil && *blacklisted {
		log.Warn(c.Request.Context(), "revoked token access attempt", domainlogger.Int64("user_id", payload.UserID))
		response.ErrorWithCode(c, http.StatusUnauthorized, string(exception.CodeInvalidToken), "token has been revoked")
		c.Abort()
		return false
	}
	c.Set(domainjwt.CtxUserIDKey, payload.UserID)
	c.Set(domainjwt.CtxTheTokenKey, parts[1])
	c.Set(domainjwt.CtxTheTokenExpireAtKey, payload.ExpiresAt)
	c.Set(domainjwt.CtxPayloadKey, payload)
	return true
}
