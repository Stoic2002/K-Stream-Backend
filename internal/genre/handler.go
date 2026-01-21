package genre

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

func (h *Handler) GetAll(c *gin.Context) {
	genres, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Failed to fetch genres", err.Error())
		return
	}
	response.Success(c, "Genres retrieved successfully", genres)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	genre, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "slug already exists" {
			response.Error(c, http.StatusConflict, err.Error(), "slug_exists")
			return
		}
		response.InternalError(c, "Failed to create genre", err.Error())
		return
	}
	response.Created(c, "Genre created successfully", genre)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	genre, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "genre not found" {
			response.NotFound(c, "Genre not found")
			return
		}
		if err.Error() == "slug already exists" {
			response.Error(c, http.StatusConflict, err.Error(), "slug_exists")
			return
		}
		response.InternalError(c, "Failed to update genre", err.Error())
		return
	}
	response.Success(c, "Genre updated successfully", genre)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete genre", err.Error())
		return
	}
	response.Success(c, "Genre deleted successfully", nil)
}
