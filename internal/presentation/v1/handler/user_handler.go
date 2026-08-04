package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type UserHandler struct {
	userService         service_contract.UserServiceInterface
	rbacService         service_contract.RBACServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
}

func NewUserHandler(
	userService service_contract.UserServiceInterface,
	rbacService service_contract.RBACServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
) *UserHandler {
	return &UserHandler{
		userService:         userService,
		rbacService:         rbacService,
		contexttokenService: contexttokenService,
	}
}

// GET /users/:id/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var hideOpt dto.HideFieldUserOption = dto.HideFieldUserOption{WalletBalance: false}

	user, err := h.userService.GetUserByID(c.Request.Context(), &tokenData, id, &hideOpt)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetProfile failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			response.Error(c, http.StatusForbidden, "You do not have permission to access this profile", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}

// PUT /users/:id/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.UpdateUserProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	req.UserID = id

	updatedUser, err := h.userService.UpdateProfile(c.Request.Context(), &tokenData, req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateProfile failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrUserAlreadyExists):
			response.Error(c, http.StatusConflict, "User with this email or phone already exists", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			response.Error(c, http.StatusForbidden, "You do not have permission to update this profile", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", updatedUser)
}

// PATCH /users/:id/status
func (h *UserHandler) ChangeStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.ChangeUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedUser, err := h.userService.ChangeStatus(c.Request.Context(), &tokenData, id, req.Status)
	if err != nil {
		log.Printf("[SERVICE ERROR] ChangeStatus failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrUserAlreadyExists):
			response.Error(c, http.StatusConflict, "User conflict occurred", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			response.Error(c, http.StatusForbidden, "You do not have permission to change status for this user", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "User status updated successfully", updatedUser)
}

// PUT /users/:id/roles
func (h *UserHandler) UpdateUserRoles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.UpdateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	req.UserID = id

	updatedUser, err := h.userService.UpdateUserRoles(c.Request.Context(), &tokenData, req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateUserRoles failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrUserAlreadyExists):
			response.Error(c, http.StatusConflict, "Role mapping conflict occurred", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			response.Error(c, http.StatusForbidden, "You do not have permission to update roles for this user", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "User roles updated successfully", updatedUser)
}

// DELETE /users/:id
func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	if err := h.userService.DeleteUserByID(c.Request.Context(), &tokenData, id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteUserByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			response.Error(c, http.StatusForbidden, "You do not have permission to delete this user", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

// GET /users?phone=...&email=...&first_name=...&last_name=...&status=...
func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var query dto.GetUsersQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters format", err)
		return
	}

	query.WalletBalance = true

	users, err := h.userService.GetUsers(ctx, &tokenData, &query)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetUsers failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidFilter):
			response.Error(c, http.StatusBadRequest, "At least one search query parameter is required (phone, email, first_name, last_name, status, permission, role)", err)

		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "No users matched the search criteria", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			response.Error(c, http.StatusForbidden, "You do not have permission to perform this search", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}
