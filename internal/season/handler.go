package season

import (
	"drakor-backend/pkg/response"
	"drakor-backend/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByDramaID(c *gin.Context) {
	dramaID := c.Param("id")
	seasons, err := h.service.GetByDramaID(c.Request.Context(), dramaID)
	if err != nil {
		response.InternalError(c, "Failed to fetch seasons", err.Error())
		return
	}
	response.Success(c, "Seasons retrieved", seasons)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	season, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to fetch season", err.Error())
		return
	}
	if season == nil {
		response.NotFound(c, "Season not found")
		return
	}
	response.Success(c, "Season detail", season)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	season, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, "Failed to create season", err.Error())
		return
	}
	response.Created(c, "Season created successfully", season)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	season, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "season not found" {
			response.NotFound(c, "Season not found")
			return
		}
		response.InternalError(c, "Failed to update season", err.Error())
		return
	}
	response.Success(c, "Season updated successfully", season)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete season", err.Error())
		return
	}
	response.Success(c, "Season deleted successfully", nil)
}
