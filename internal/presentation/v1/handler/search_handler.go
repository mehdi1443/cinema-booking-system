package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type SearchHandler struct {
	searchService service_contract.SearchServiceInterface
	log           domainlogger.Logger
}

func NewSearchHandler(
	searchService service_contract.SearchServiceInterface,
	log domainlogger.Logger,
) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		log:           log,
	}
}

// -----------------------------------------------------------------------------
// GET /search/suggest
// -----------------------------------------------------------------------------

func (h *SearchHandler) Suggest(c *gin.Context) {
	var req dto.SuggestRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		HandleBadRequest(c, "Invalid request parameters")
		return
	}

	results, err := h.searchService.Suggest(c.Request.Context(), req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Search suggestions retrieved successfully", results)
}

// -----------------------------------------------------------------------------
// GET /search
// -----------------------------------------------------------------------------

func (h *SearchHandler) Search(c *gin.Context) {
	var req dto.SearchQueryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		HandleBadRequest(c, "Invalid request parameters")
		return
	}

	results, total, err := h.searchService.Search(c.Request.Context(), req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	meta := response.Meta{
		Page:  1,
		Limit: req.Limit,
		Total: int(total),
	}

	response.SuccessWithMeta(c, http.StatusOK, "Search results retrieved successfully", results, meta)
}

func (h *SearchHandler) SyncAll(c *gin.Context) {
	if err := h.searchService.SyncAll(c.Request.Context()); err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "All entities synced to Elasticsearch successfully", nil)
}
