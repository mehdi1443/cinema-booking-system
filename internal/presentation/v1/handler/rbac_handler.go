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

type RBACHandler struct {
	rbacService         service_contract.RBACServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
}

func NewRBACHandler(
	rbacService service_contract.RBACServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
) *RBACHandler {
	return &RBACHandler{
		rbacService:         rbacService,
		contexttokenService: contexttokenService,
	}
}

// GetUserPermissions handles GET /users/:id/permissions
func (h *RBACHandler) GetUserPermissions(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

	permissions, err := h.rbacService.GetUserPermissions(c.Request.Context(), &tokenData, userID)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetUserPermissions failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to retrieve user permissions", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "User permissions retrieved successfully", permissions)
}

// GetUserRoles handles GET /users/:id/roles
func (h *RBACHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

	roles, err := h.rbacService.GetUserRoles(c.Request.Context(), &tokenData, userID)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetUserRoles failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to retrieve user roles", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "User roles retrieved successfully", roles)
}

// AssignPermissionsToRole handles POST /roles/:id/permissions
func (h *RBACHandler) AssignPermissionsToRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.AssignPermissionsToRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.rbacService.AssignPermissionsToRole(c.Request.Context(), &tokenData, roleID, req.PermissionIDs); err != nil {
		log.Printf("[SERVICE ERROR] AssignPermissionsToRole failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrRoleNotFound):
			response.Error(c, http.StatusNotFound, "Role not found", err)
		case errors.Is(err, exception.ErrPermissionNotFound):
			response.Error(c, http.StatusNotFound, "One or more permissions not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to assign permissions to role", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Permissions assigned to role successfully", nil)
}

// DischargePermissionsOfRole handles DELETE /roles/:id/permissions
func (h *RBACHandler) DischargePermissionsOfRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.DischargePermissionsOfRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.rbacService.DischargePermissionsOfRole(c.Request.Context(), &tokenData, roleID, req.PermissionIDs); err != nil {
		log.Printf("[SERVICE ERROR] DischargePermissionsOfRole failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrRoleNotFound):
			response.Error(c, http.StatusNotFound, "Role not found", err)
		case errors.Is(err, exception.ErrPermissionNotFound):
			response.Error(c, http.StatusNotFound, "One or more permissions not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to discharge permissions from role", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Permissions discharged from role successfully", nil)
}

// GetRoleByID handles GET /roles/:id
func (h *RBACHandler) GetRoleByID(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	role, err := h.rbacService.GetRoleByID(c.Request.Context(), &tokenData, roleID)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetRoleByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrRoleNotFound):
			response.Error(c, http.StatusNotFound, "Role not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Role retrieved successfully", role)
}

// GetPermissionByID handles GET /permissions/:id
func (h *RBACHandler) GetPermissionByID(c *gin.Context) {
	permissionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid permission ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	permission, err := h.rbacService.GetPermissionsByID(c.Request.Context(), &tokenData, permissionID)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetPermissionByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrPermissionNotFound):
			response.Error(c, http.StatusNotFound, "Permission not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Permission retrieved successfully", permission)
}

// GetPermissionsByRoleName handles GET /permissions/
func (h *RBACHandler) GetPermissionsByRoleName(c *gin.Context) {
	roleName := c.Param("name")
	if roleName == "" {
		response.Error(c, http.StatusBadRequest, "Role name is required", nil)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	permissions, err := h.rbacService.GetPermissionsByRoleName(c.Request.Context(), &tokenData, roleName)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetPermissionsByRoleName failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrRoleNotFound):
			response.Error(c, http.StatusNotFound, "Role not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to fetch permissions for role", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Role permissions retrieved successfully", permissions)
}
