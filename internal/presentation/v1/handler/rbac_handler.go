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

type RBACHandler struct {
	rbacService         service_contract.RBACServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 logger.Logger
}

func NewRBACHandler(
	rbacService service_contract.RBACServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log logger.Logger,
) *RBACHandler {
	return &RBACHandler{
		rbacService:         rbacService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// AppendPermissionsToRole handles POST /roles/:id/permissions
func (h *RBACHandler) AppendPermissionsToRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid role ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "AppendPermissionsToRole missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.AppendPermissionsToRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for AppendPermissionsToRole", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.rbacService.AppendPermissionsToRole(c.Request.Context(), &tokenData, roleID, req.PermissionIDs); err != nil {
		switch {
		case errors.Is(err, exception.ErrRoleNotFound):
			h.log.Warn(logger.Validation, logger.Api, "AppendPermissionsToRole failed: role not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Role not found", err)
		case errors.Is(err, exception.ErrPermissionNotFound):
			h.log.Warn(logger.Validation, logger.Api, "AppendPermissionsToRole failed: permission not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "One or more permissions not found", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "AppendPermissionsToRole internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Failed to append permissions to role", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Permissions appended to role successfully", nil)
}

// DeletePermissionsOfRole handles DELETE /roles/:id/permissions
func (h *RBACHandler) DeletePermissionsOfRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid role ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "DeletePermissionsOfRole missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.DeletePermissionsOfRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for DeletePermissionsOfRole", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.rbacService.DeletePermissionsOfRole(c.Request.Context(), &tokenData, roleID, req.PermissionIDs); err != nil {
		switch {
		case errors.Is(err, exception.ErrRoleNotFound):
			h.log.Warn(logger.Validation, logger.Api, "DeletePermissionsOfRole failed: role not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Role not found", err)
		case errors.Is(err, exception.ErrPermissionNotFound):
			h.log.Warn(logger.Validation, logger.Api, "DeletePermissionsOfRole failed: permission not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "One or more permissions not found", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "DeletePermissionsOfRole internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Failed to delete permissions from role", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Permissions deleted from role successfully", nil)
}

// GetRoleByID handles GET /roles/:id
func (h *RBACHandler) GetRoleByID(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid role ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "GetRoleByID missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	role, err := h.rbacService.GetRoleByID(c.Request.Context(), &tokenData, roleID)
	if err != nil {
		if errors.Is(err, exception.ErrRoleNotFound) {
			h.log.Warn(logger.Validation, logger.Api, "GetRoleByID: role not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Role not found", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "GetRoleByID internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		return
	}

	response.Success(c, http.StatusOK, "Role retrieved successfully", role)
}

// GetPermissionByID handles GET /permissions/:id
func (h *RBACHandler) GetPermissionByID(c *gin.Context) {
	permissionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid permission ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid permission ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "GetPermissionByID missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	permission, err := h.rbacService.GetPermissionsByID(c.Request.Context(), &tokenData, permissionID)
	if err != nil {
		if errors.Is(err, exception.ErrPermissionNotFound) {
			h.log.Warn(logger.Validation, logger.Api, "GetPermissionByID: permission not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Permission not found", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "GetPermissionByID internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		return
	}

	response.Success(c, http.StatusOK, "Permission retrieved successfully", permission)
}

// GetPermissionsByRoleName handles GET /permissions/
func (h *RBACHandler) GetPermissionsByRoleName(c *gin.Context) {
	roleName := c.Param("name")
	if roleName == "" {
		h.log.Warn(logger.Validation, logger.Api, "Role name empty in GetPermissionsByRoleName", nil)
		response.Error(c, http.StatusBadRequest, "Role name is required", nil)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "GetPermissionsByRoleName missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	permissions, err := h.rbacService.GetPermissionsByRoleName(c.Request.Context(), &tokenData, roleName)
	if err != nil {
		if errors.Is(err, exception.ErrRoleNotFound) {
			h.log.Warn(logger.Validation, logger.Api, "GetPermissionsByRoleName: role not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Role not found", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "GetPermissionsByRoleName internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to fetch permissions for role", err)
		return
	}

	response.Success(c, http.StatusOK, "Role permissions retrieved successfully", permissions)
}

// CreateRole handles POST /roles
func (h *RBACHandler) CreateRole(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "CreateRole missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for CreateRole", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	role, err := h.rbacService.CreateRole(c.Request.Context(), &tokenData, &req)
	if err != nil {
		if errors.Is(err, exception.ErrRoleAlreadyExists) {
			h.log.Warn(logger.Validation, logger.Api, "CreateRole failed: role already exists", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusConflict, "Role already exists", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "CreateRole internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to create role", err)
		return
	}

	response.Success(c, http.StatusCreated, "Role created successfully", role)
}

// UpdateRole handles PUT /roles/:id
func (h *RBACHandler) UpdateRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid role ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "UpdateRole missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for UpdateRole", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	role, err := h.rbacService.UpdateRole(c.Request.Context(), &tokenData, roleID, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrRoleNotFound):
			h.log.Warn(logger.Validation, logger.Api, "UpdateRole failed: role not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Role not found", err)
		case errors.Is(err, exception.ErrRoleAlreadyExists):
			h.log.Warn(logger.Validation, logger.Api, "UpdateRole failed: role name in use", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusConflict, "Role name already in use", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "UpdateRole internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Failed to update role", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Role updated successfully", role)
}
