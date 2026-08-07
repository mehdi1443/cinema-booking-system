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

type WalletHandler struct {
	walletService       service_contract.WalletServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 logger.Logger
}

func NewWalletHandler(
	walletService service_contract.WalletServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log logger.Logger,
) *WalletHandler {
	return &WalletHandler{
		walletService:       walletService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// GetBalance handles GET /users/:id/wallet/balance
func (h *WalletHandler) GetBalance(c *gin.Context) {
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
		h.log.Warn(logger.Internal, logger.Api, "GetBalance missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	balance, err := h.walletService.GetBalance(c.Request.Context(), &tokenData, userID)
	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			h.log.Warn(logger.Validation, logger.Api, "GetBalance: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "GetBalance internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve balance", err)
		return
	}

	response.Success(c, http.StatusOK, "Balance retrieved successfully", dto.WalletBalanceResponse{
		UserID:  userID,
		Balance: balance,
	})
}

// ChangeBalance handles POST /users/:id/wallet/balance
func (h *WalletHandler) ChangeBalance(c *gin.Context) {
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
		h.log.Warn(logger.Internal, logger.Api, "ChangeBalance missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.ChangeBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for ChangeBalance", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	newBalance, err := h.walletService.ChangeBalance(c.Request.Context(), &tokenData, userID, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "ChangeBalance failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)
		case errors.Is(err, exception.ErrInvalidAmount):
			h.log.Warn(logger.Validation, logger.Api, "ChangeBalance failed: invalid amount", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid amount provided", err)
		case errors.Is(err, exception.ErrInsufficientBalance):
			h.log.Warn(logger.Validation, logger.Api, "ChangeBalance failed: insufficient balance", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Insufficient wallet balance", err)
		case errors.Is(err, exception.ErrReachManBalance):
			h.log.Warn(logger.Validation, logger.Api, "ChangeBalance failed: max wallet limit reached", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Maximum wallet balance limit reached", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "ChangeBalance internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Failed to update balance", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Balance increased successfully", dto.WalletBalanceResponse{
		UserID:  userID,
		Balance: newBalance,
	})
}
