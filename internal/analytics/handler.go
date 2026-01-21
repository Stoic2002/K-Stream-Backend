package analytics

import (
	"drakor-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetDashboard(c *gin.Context) {
	stats, err := h.service.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Failed to fetch dashboard stats", err.Error())
		return
	}
	response.Success(c, "Dashboard stats retrieved", stats)
}
