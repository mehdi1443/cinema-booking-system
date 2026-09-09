package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
)

func parseParamID(c *gin.Context, paramName string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(paramName), 10, 64)
	if err != nil || id <= 0 {
		HandleBadRequest(c, "Invalid "+paramName+" parameter")
		return 0, false
	}
	return id, true
}

func hasPermission(c *gin.Context, required enum.PermissionType) bool {
	value, exists := c.Get(domainjwt.CtxPermissionsKey)
	if !exists {
		return false
	}
	permissions, ok := value.([]string)
	if !ok {
		return false
	}
	for _, permission := range permissions {
		if permission == required.String() || permission == enum.PermissionAll.String() {
			return true
		}
	}
	return false
}

func getAuthUserID(c *gin.Context, log domainlogger.Logger) (int64, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		Handle(c, exception.ErrGetTokenFromContext, log)
		return 0, false
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		Handle(c, exception.ErrGetTokenFromContext, log)
		return 0, false
	}

	return userID, true
}

func getOptionalAuthUserID(c *gin.Context) *int64 {
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(int64); ok {
			return &uid
		}
	}
	return nil
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return page, limit
}
