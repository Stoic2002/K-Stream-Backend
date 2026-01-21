package episode

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

func (h *Handler) GetBySeasonID(c *gin.Context) {
	seasonID := c.Param("id")
	episodes, err := h.service.GetBySeasonID(c.Request.Context(), seasonID)
	if err != nil {
		response.InternalError(c, "Failed to fetch episodes", err.Error())
		return
	}
	response.Success(c, "Episodes retrieved", episodes)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	episode, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to fetch episode", err.Error())
		return
	}
	if episode == nil {
		response.NotFound(c, "Episode not found")
		return
	}
	response.Success(c, "Episode detail", episode)
}

func (h *Handler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req CreateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	episode, err := h.service.Create(c.Request.Context(), userID.(string), req)
	if err != nil {
		response.InternalError(c, "Failed to create episode", err.Error())
		return
	}
	response.Created(c, "Episode created successfully", episode)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	episode, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "episode not found" {
			response.NotFound(c, "Episode not found")
			return
		}
		response.InternalError(c, "Failed to update episode", err.Error())
		return
	}
	response.Success(c, "Episode updated successfully", episode)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete episode", err.Error())
		return
	}
	response.Success(c, "Episode deleted successfully", nil)
}
