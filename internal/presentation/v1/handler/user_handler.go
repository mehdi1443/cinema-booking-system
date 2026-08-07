package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type UserHandler struct {
	userService         service_contract.UserServiceInterface
	rbacService         service_contract.RBACServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 logger.Logger
}

func NewUserHandler(
	userService service_contract.UserServiceInterface,
	rbacService service_contract.RBACServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log logger.Logger,
) *UserHandler {
	return &UserHandler{
		userService:         userService,
		rbacService:         rbacService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// TODO : define and decide userStatus BussiznesLogic
// GET /users/:id/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid user ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "GetProfile missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var hideOpt dto.HideFieldUserOption = dto.HideFieldUserOption{WalletBalance: false}

	user, err := h.userService.GetUserByID(c.Request.Context(), &tokenData, id, &hideOpt)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "GetProfile failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			h.log.Warn(logger.Validation, logger.Api, "GetProfile failed: unauthorized access", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusForbidden, "You do not have permission to access this profile", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Select, "GetProfile database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "GetProfile unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Validation, logger.Api, "Invalid user ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "UpdateProfile missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.UpdateUserProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for UpdateProfile", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedUser, err := h.userService.UpdateProfile(c.Request.Context(), &tokenData, id, req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "UpdateProfile failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrUserAlreadyExists):
			h.log.Warn(logger.Validation, logger.Api, "UpdateProfile failed: conflict on email/phone", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusConflict, "User with this email or phone already exists", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			h.log.Warn(logger.Validation, logger.Api, "UpdateProfile failed: unauthorized access", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusForbidden, "You do not have permission to update this profile", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Update, "UpdateProfile database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "UpdateProfile unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Validation, logger.Api, "Invalid user ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "ChangeStatus missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.ChangeUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for ChangeStatus", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedUser, err := h.userService.ChangeStatus(c.Request.Context(), &tokenData, id, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "ChangeStatus failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrUserAlreadyExists):
			h.log.Warn(logger.Validation, logger.Api, "ChangeStatus failed: conflict", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusConflict, "User conflict occurred", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			h.log.Warn(logger.Validation, logger.Api, "ChangeStatus failed: forbidden", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusForbidden, "You do not have permission to change status for this user", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Update, "ChangeStatus database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "ChangeStatus unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Validation, logger.Api, "Invalid user ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "UpdateUserRoles missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.UpdateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for UpdateUserRoles", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedUser, err := h.rbacService.UpdateUserRoles(c.Request.Context(), &tokenData, id, req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "UpdateUserRoles failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrUserAlreadyExists):
			h.log.Warn(logger.Validation, logger.Api, "UpdateUserRoles failed: role mapping conflict", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusConflict, "Role mapping conflict occurred", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			h.log.Warn(logger.Validation, logger.Api, "UpdateUserRoles failed: forbidden", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusForbidden, "You do not have permission to update roles for this user", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Update, "UpdateUserRoles database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "UpdateUserRoles unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Validation, logger.Api, "Invalid user ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "DeleteUserByID missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	if err := h.userService.DeleteUserByID(c.Request.Context(), &tokenData, id); err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "DeleteUserByID failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			h.log.Warn(logger.Validation, logger.Api, "DeleteUserByID failed: forbidden", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusForbidden, "You do not have permission to delete this user", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Delete, "DeleteUserByID database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "DeleteUserByID unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Internal, logger.Api, "GetUsers missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var query dto.GetUsersQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid query parameters for GetUsers", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid query parameters format", err)
		return
	}

	query.WalletBalance = true

	users, err := h.userService.GetUsers(ctx, &tokenData, &query)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrInvalidFilter):
			h.log.Warn(logger.Validation, logger.Api, "GetUsers failed: at least one filter required", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "At least one search query parameter is required (phone, email, first_name, last_name, status, permission, role)", err)

		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "GetUsers: no users matched criteria", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "No users matched the search criteria", err)

		case errors.Is(err, exception.ErrNoteAccessible):
			h.log.Warn(logger.Validation, logger.Api, "GetUsers failed: forbidden", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusForbidden, "You do not have permission to perform this search", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Select, "GetUsers database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "GetUsers unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid user ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "UpdatePassword missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for UpdatePassword", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	success, err := h.userService.UpdatePassword(c.Request.Context(), &tokenData, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "UpdatePassword failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)
		case errors.Is(err, exception.ErrInvalidPassword):
			h.log.Warn(logger.Validation, logger.Api, "UpdatePassword failed: invalid password", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid password provided", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "UpdatePassword internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Failed to update password", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Password updated successfully", gin.H{
		"success": success,
	})
}

// GetUserPermissions handles GET /users/:id/permissions
func (h *UserHandler) GetUserPermissions(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid user ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "GetUserPermissions missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	permissions, err := h.rbacService.GetUserPermissions(c.Request.Context(), &tokenData, userID)
	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			h.log.Warn(logger.Validation, logger.Api, "GetUserPermissions: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "GetUserPermissions internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve user permissions", err)
		return
	}

	response.Success(c, http.StatusOK, "User permissions retrieved successfully", permissions)
}

// GetUserRoles handles GET /users/:id/roles
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid user ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "GetUserRoles missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	roles, err := h.rbacService.GetUserRoles(c.Request.Context(), &tokenData, userID)
	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			h.log.Warn(logger.Validation, logger.Api, "GetUserRoles: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "GetUserRoles internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve user roles", err)
		return
	}

	response.Success(c, http.StatusOK, "User roles retrieved successfully", roles)
}
