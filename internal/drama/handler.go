package drama

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

func (h *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	query := c.Query("q")
	genreID := c.Query("genre")
	status := c.Query("status")
	sort := c.Query("sort")

	dramas, total, err := h.service.GetAll(c.Request.Context(), page, limit, query, genreID, status, sort)
	if err != nil {
		response.InternalError(c, "Failed to fetch dramas", err.Error())
		return
	}
	response.Paginated(c, dramas, total, page, limit)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	drama, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to fetch drama", err.Error())
		return
	}
	if drama == nil {
		response.NotFound(c, "Drama not found")
		return
	}
	response.Success(c, "Drama detail", drama)
}

func (h *Handler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req CreateDramaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	drama, err := h.service.Create(c.Request.Context(), userID.(string), req)
	if err != nil {
		response.InternalError(c, "Failed to create drama", err.Error())
		return
	}
	response.Created(c, "Drama created successfully", drama)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateDramaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	drama, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "drama not found" {
			response.NotFound(c, "Drama not found")
			return
		}
		response.InternalError(c, "Failed to update drama", err.Error())
		return
	}
	response.Success(c, "Drama updated successfully", drama)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete drama", err.Error())
		return
	}
	response.Success(c, "Drama deleted successfully", nil)
}
