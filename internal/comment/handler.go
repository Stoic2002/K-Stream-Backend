package comment

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

func (h *Handler) GetByEpisode(c *gin.Context) {
	episodeID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	comments, total, err := h.service.GetByEpisode(c.Request.Context(), episodeID, page, limit)
	if err != nil {
		response.InternalError(c, "Failed to fetch comments", err.Error())
		return
	}
	response.Paginated(c, comments, total, page, limit)
}

func (h *Handler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	comment, err := h.service.Create(c.Request.Context(), userID.(string), req)
	if err != nil {
		response.InternalError(c, "Failed to create comment", err.Error())
		return
	}
	response.Created(c, "Comment created successfully", comment)
}

func (h *Handler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	commentID := c.Param("id")

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	comment, err := h.service.Update(c.Request.Context(), userID.(string), commentID, req)
	if err != nil {
		response.InternalError(c, "Failed to update comment", err.Error())
		return
	}
	response.Success(c, "Comment updated successfully", comment)
}

func (h *Handler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userRole, _ := c.Get("userRole")
	isAdmin := userRole == "admin"

	commentID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), userID.(string), commentID, isAdmin); err != nil {
		response.InternalError(c, "Failed to delete comment", err.Error())
		return
	}
	response.Success(c, "Comment deleted successfully", nil)
}
