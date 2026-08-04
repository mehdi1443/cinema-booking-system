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

type WalletHandler struct {
	walletService       service_contract.WalletServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
}

func NewWalletHandler(
	walletService service_contract.WalletServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
) *WalletHandler {
	return &WalletHandler{
		walletService:       walletService,
		contexttokenService: contexttokenService,
	}
}

// GetBalance handles GET /users/:id/wallet/balance
func (h *WalletHandler) GetBalance(c *gin.Context) {
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

	balance, err := h.walletService.GetBalance(c.Request.Context(), &tokenData, userID)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetBalance failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to retrieve balance", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Balance retrieved successfully", dto.WalletBalanceResponse{
		UserID:  userID,
		Balance: balance,
	})
}

// IncreaseBalance handles POST /users/:id/wallet/balance
func (h *WalletHandler) ChangeBalance(c *gin.Context) {
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

	var req dto.ChangeBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	newBalance, err := h.walletService.ChangeBalance(c.Request.Context(), &tokenData, userID, req.Amount)
	if err != nil {
		log.Printf("[SERVICE ERROR] IncreaseBalance failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)
		case errors.Is(err, exception.ErrInvalidAmount):
			response.Error(c, http.StatusBadRequest, "Invalid amount provided", err)
		case errors.Is(err, exception.ErrInsufficientBalance):
			response.Error(c, http.StatusBadRequest, "Insufficient wallet balance", err)
		case errors.Is(err, exception.ErrReachManBalance):
			response.Error(c, http.StatusBadRequest, "Maximum wallet balance limit reached", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to update balance", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Balance increased successfully", dto.WalletBalanceResponse{
		UserID:  userID,
		Balance: newBalance,
	})
}
