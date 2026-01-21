package history

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

func (h *Handler) Record(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req RecordHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	if err := h.service.RecordProgress(c.Request.Context(), userID.(string), req.EpisodeID, req.ProgressSeconds, req.IsFinished); err != nil {
		response.InternalError(c, "Failed to record progress", err.Error())
		return
	}
	response.Success(c, "Progress recorded", nil)
}

func (h *Handler) GetMine(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.service.GetMyHistory(c.Request.Context(), userID.(string), page, limit)
	if err != nil {
		response.InternalError(c, "Failed to fetch history", err.Error())
		return
	}
	response.Paginated(c, items, total, page, limit)
}

func (h *Handler) GetProgress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	episodeID := c.Param("episodeID")

	history, err := h.service.GetEpisodeProgress(c.Request.Context(), userID.(string), episodeID)
	if err != nil {
		response.InternalError(c, "Failed to fetch progress", err.Error())
		return
	}
	// Note: It's okay to return null history if not watched yet.
	response.Success(c, "Progress retrieved", history)
}
