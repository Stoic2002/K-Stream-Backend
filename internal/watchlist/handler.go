package watchlist

import (
	"drakor-backend/pkg/response"
	"drakor-backend/pkg/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Add(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req AddWatchlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	if err := h.service.AddToWatchlist(c.Request.Context(), userID.(string), req.DramaID); err != nil {
		response.InternalError(c, "Failed to add to watchlist", err.Error())
		return
	}
	response.Success(c, "Added to watchlist", nil)
}

func (h *Handler) Remove(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	dramaID := c.Param("dramaID")

	if err := h.service.RemoveFromWatchlist(c.Request.Context(), userID.(string), dramaID); err != nil {
		response.InternalError(c, "Failed to remove from watchlist", err.Error())
		return
	}
	response.Success(c, "Removed from watchlist", nil)
}

func (h *Handler) GetMine(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.service.GetMyWatchlist(c.Request.Context(), userID.(string), page, limit)
	if err != nil {
		response.InternalError(c, "Failed to fetch watchlist", err.Error())
		return
	}
	response.Paginated(c, items, total, page, limit)
}

func (h *Handler) Check(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	dramaID := c.Param("dramaID")

	isWatchlisted, err := h.service.CheckIsWatchlisted(c.Request.Context(), userID.(string), dramaID)
	if err != nil {
		response.InternalError(c, "Failed to check status", err.Error())
		return
	}
	response.Success(c, "Check status success", gin.H{"is_watchlisted": isWatchlisted})
}
