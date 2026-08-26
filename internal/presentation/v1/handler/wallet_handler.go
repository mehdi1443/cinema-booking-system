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

type WalletHandler struct {
	walletService       service_contract.WalletServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 domainlogger.Logger
}

func NewWalletHandler(
	walletService service_contract.WalletServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log domainlogger.Logger,
) *WalletHandler {
	return &WalletHandler{
		walletService:       walletService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// -----------------------------------------------------------------------------
// GET /users/:id/wallet/balance
// -----------------------------------------------------------------------------

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	balance, err := h.walletService.GetBalance(c.Request.Context(), &tokenData, userID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Balance retrieved successfully", dto.WalletBalanceResponse{
		UserID:  userID,
		Balance: balance,
	})
}

// -----------------------------------------------------------------------------
// POST /users/:id/wallet/balance
// -----------------------------------------------------------------------------

func (h *WalletHandler) ChangeBalance(c *gin.Context) {
	userID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.ChangeBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	newBalance, err := h.walletService.ChangeBalance(c.Request.Context(), &tokenData, userID, req.Amount)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Balance updated successfully", dto.WalletBalanceResponse{
		UserID:  userID,
		Balance: newBalance,
	})
}
