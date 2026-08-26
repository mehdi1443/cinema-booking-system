package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type UserHandler struct {
	userService         service_contract.UserServiceInterface
	rbacService         service_contract.RBACServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 domainlogger.Logger
}

func NewUserHandler(
	userService service_contract.UserServiceInterface,
	rbacService service_contract.RBACServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log domainlogger.Logger,
) *UserHandler {
	return &UserHandler{
		userService:         userService,
		rbacService:         rbacService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// -----------------------------------------------------------------------------
// GET /users/:id/profile
// -----------------------------------------------------------------------------

func (h *UserHandler) GetProfile(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	hideOpt := dto.HideFieldUserOption{WalletBalance: false}

	user, err := h.userService.GetUserByID(c.Request.Context(), &tokenData, id, &hideOpt)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}

// -----------------------------------------------------------------------------
// PUT /users/:id/profile
// -----------------------------------------------------------------------------

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.UpdateUserProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedUser, err := h.userService.UpdateProfile(c.Request.Context(), &tokenData, id, req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", updatedUser)
}

// -----------------------------------------------------------------------------
// PATCH /users/:id/status
// -----------------------------------------------------------------------------

func (h *UserHandler) ChangeStatus(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.ChangeUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedUser, err := h.userService.ChangeStatus(c.Request.Context(), &tokenData, id, req.Status)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "User status updated successfully", updatedUser)
}

// -----------------------------------------------------------------------------
// PUT /users/:id/roles
// -----------------------------------------------------------------------------

func (h *UserHandler) UpdateUserRoles(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.UpdateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedUser, err := h.rbacService.UpdateUserRoles(c.Request.Context(), &tokenData, id, req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "User roles updated successfully", updatedUser)
}

// -----------------------------------------------------------------------------
// DELETE /users/:id
// -----------------------------------------------------------------------------

func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	if err := h.userService.DeleteUserByID(c.Request.Context(), &tokenData, id); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

// -----------------------------------------------------------------------------
// GET /users
// -----------------------------------------------------------------------------

func (h *UserHandler) GetUsers(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var query dto.GetUsersQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleBadRequest(c, "Invalid query parameters format")
		return
	}

	query.WalletBalance = true

	users, err := h.userService.GetUsers(c.Request.Context(), &tokenData, &query)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}

// -----------------------------------------------------------------------------
// PUT /users/:id/password
// -----------------------------------------------------------------------------

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	success, err := h.userService.UpdatePassword(c.Request.Context(), &tokenData, userID, req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Password updated successfully", gin.H{
		"success": success,
	})
}

// -----------------------------------------------------------------------------
// GET /users/:id/permissions
// -----------------------------------------------------------------------------

func (h *UserHandler) GetUserPermissions(c *gin.Context) {
	userID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	permissions, err := h.rbacService.GetUserPermissions(c.Request.Context(), &tokenData, userID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "User permissions retrieved successfully", permissions)
}

// -----------------------------------------------------------------------------
// GET /users/:id/roles
// -----------------------------------------------------------------------------

func (h *UserHandler) GetUserRoles(c *gin.Context) {
	userID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	roles, err := h.rbacService.GetUserRoles(c.Request.Context(), &tokenData, userID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "User roles retrieved successfully", roles)
}
