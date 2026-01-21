package auth

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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input data", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "email already registered" {
			response.Error(c, http.StatusConflict, err.Error(), "email_exists")
			return
		}
		response.InternalError(c, "Failed to register user", err.Error())
		return
	}

	response.Created(c, "User registered successfully", resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input data", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to login", err.Error())
		return
	}

	response.Success(c, "Login successful", resp)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		response.InternalError(c, "Failed to get profile", err.Error())
		return
	}

	response.Success(c, "User profile retrieved", user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input data", err.Error())
		return
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		response.Error(c, http.StatusBadRequest, "Validation failed", "validation_error")
		return
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), userID.(string), req)
	if err != nil {
		response.InternalError(c, "Failed to update profile", err.Error())
		return
	}

	response.Success(c, "Profile updated successfully", user)
}

// --- Admin Handlers ---

func (h *Handler) GetAllUsers(c *gin.Context) {
	page := 1 // Default
	limit := 10

	users, total, err := h.service.GetAllUsers(c.Request.Context(), page, limit)
	if err != nil {
		response.InternalError(c, "Failed to fetch users", err.Error())
		return
	}
	response.Paginated(c, users, total, page, limit)
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		Role string `json:"role" binding:"required,oneof=user admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input", err.Error())
		return
	}

	if err := h.service.UpdateUserRole(c.Request.Context(), userID, req.Role); err != nil {
		response.InternalError(c, "Failed to update role", err.Error())
		return
	}
	response.Success(c, "User role updated successfully", nil)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.DeleteUser(c.Request.Context(), userID); err != nil {
		response.InternalError(c, "Failed to delete user", err.Error())
		return
	}
	response.Success(c, "User deleted successfully", nil)
}
