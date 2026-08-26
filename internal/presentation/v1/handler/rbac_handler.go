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

type RBACHandler struct {
	rbacService         service_contract.RBACServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 domainlogger.Logger
}

func NewRBACHandler(
	rbacService service_contract.RBACServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log domainlogger.Logger,
) *RBACHandler {
	return &RBACHandler{
		rbacService:         rbacService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// -----------------------------------------------------------------------------
// POST /roles/:id/permissions
// -----------------------------------------------------------------------------

func (h *RBACHandler) AppendPermissionsToRole(c *gin.Context) {
	roleID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.AppendPermissionsToRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	if err := h.rbacService.AppendPermissionsToRole(c.Request.Context(), &tokenData, roleID, req.PermissionIDs); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Permissions appended to role successfully", nil)
}

// -----------------------------------------------------------------------------
// DELETE /roles/:id/permissions
// -----------------------------------------------------------------------------

func (h *RBACHandler) DeletePermissionsOfRole(c *gin.Context) {
	roleID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.DeletePermissionsOfRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	if err := h.rbacService.DeletePermissionsOfRole(c.Request.Context(), &tokenData, roleID, req.PermissionIDs); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Permissions deleted from role successfully", nil)
}

// -----------------------------------------------------------------------------
// GET /roles/:id
// -----------------------------------------------------------------------------

func (h *RBACHandler) GetRoleByID(c *gin.Context) {
	roleID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	role, err := h.rbacService.GetRoleByID(c.Request.Context(), &tokenData, roleID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Role retrieved successfully", role)
}

// -----------------------------------------------------------------------------
// GET /permissions/:id
// -----------------------------------------------------------------------------

func (h *RBACHandler) GetPermissionByID(c *gin.Context) {
	permissionID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	permission, err := h.rbacService.GetPermissionsByID(c.Request.Context(), &tokenData, permissionID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Permission retrieved successfully", permission)
}

// -----------------------------------------------------------------------------
// GET /permissions/by-role/:name
// -----------------------------------------------------------------------------

func (h *RBACHandler) GetPermissionsByRoleName(c *gin.Context) {
	roleName := c.Param("name")
	if roleName == "" {
		HandleBadRequest(c, "Role name is required")
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	permissions, err := h.rbacService.GetPermissionsByRoleName(c.Request.Context(), &tokenData, roleName)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Role permissions retrieved successfully", permissions)
}

// -----------------------------------------------------------------------------
// POST /roles
// -----------------------------------------------------------------------------

func (h *RBACHandler) CreateRole(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	role, err := h.rbacService.CreateRole(c.Request.Context(), &tokenData, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusCreated, "Role created successfully", role)
}

// -----------------------------------------------------------------------------
// PUT /roles/:id
// -----------------------------------------------------------------------------

func (h *RBACHandler) UpdateRole(c *gin.Context) {
	roleID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	role, err := h.rbacService.UpdateRole(c.Request.Context(), &tokenData, roleID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Role updated successfully", role)
}
