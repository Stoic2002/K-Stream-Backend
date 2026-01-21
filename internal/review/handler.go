package review

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

func (h *Handler) GetByDrama(c *gin.Context) {
	dramaID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.service.GetByDrama(c.Request.Context(), dramaID, page, limit)
	if err != nil {
		response.InternalError(c, "Failed to fetch reviews", err.Error())
		return
	}
	response.Paginated(c, reviews, total, page, limit)
}

func (h *Handler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	review, err := h.service.Create(c.Request.Context(), userID.(string), req)
	if err != nil {
		if err.Error() == "review already exists" {
			response.Error(c, http.StatusConflict, "You have already reviewed this drama", "conflict_error")
			return
		}
		response.InternalError(c, "Failed to create review", err.Error())
		return
	}
	response.Created(c, "Review created successfully", review)
}

func (h *Handler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	reviewID := c.Param("id")

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	review, err := h.service.Update(c.Request.Context(), userID.(string), reviewID, req)
	if err != nil {
		response.InternalError(c, "Failed to update review", err.Error())
		return
	}
	response.Success(c, "Review updated successfully", review)
}

func (h *Handler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userRole, _ := c.Get("userRole")
	isAdmin := userRole == "admin"

	reviewID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), userID.(string), reviewID, isAdmin); err != nil {
		response.InternalError(c, "Failed to delete review", err.Error())
		return
	}
	response.Success(c, "Review deleted successfully", nil)
}
