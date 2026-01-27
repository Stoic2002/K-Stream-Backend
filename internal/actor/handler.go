package actor

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
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	actors, total, err := h.service.GetAll(c.Request.Context(), page, limit, search)
	if err != nil {
		response.InternalError(c, "Failed to fetch actors", err.Error())
		return
	}
	response.Paginated(c, actors, total, page, limit)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	actor, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to fetch actor", err.Error())
		return
	}
	if actor == nil {
		response.NotFound(c, "Actor not found")
		return
	}
	response.Success(c, "Actor detail", actor)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateActorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	actor, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, "Failed to create actor", err.Error())
		return
	}
	response.Created(c, "Actor created successfully", actor)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateActorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	actor, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "actor not found" {
			response.NotFound(c, "Actor not found")
			return
		}
		response.InternalError(c, "Failed to update actor", err.Error())
		return
	}
	response.Success(c, "Actor updated successfully", actor)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete actor", err.Error())
		return
	}
	response.Success(c, "Actor deleted successfully", nil)
}
