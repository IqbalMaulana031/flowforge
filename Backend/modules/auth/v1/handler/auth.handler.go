package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"flowforge-api/middleware"
	"flowforge-api/modules/auth/v1/service"
	"flowforge-api/resource"
	"flowforge-api/response"
)

type AuthHandler struct{ service service.AuthUseCase }

func NewAuthHandler(service service.AuthUseCase) *AuthHandler { return &AuthHandler{service: service} }
func (h *AuthHandler) Register(c *gin.Context) {
	var req resource.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "REGISTER_FAILED", err.Error())
		return
	}
	response.Created(c, res)
}
func (h *AuthHandler) Login(c *gin.Context) {
	var req resource.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "LOGIN_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *AuthHandler) Profile(c *gin.Context) {
	res, err := h.service.Profile(c.Request.Context(), c.GetString(middleware.UserIDKey))
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req resource.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.UpdateProfile(c.Request.Context(), c.GetString(middleware.UserIDKey), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req resource.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "REFRESH_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.service.Logout(c.Request.Context(), c.GetString(middleware.AccessTokenKey)); err != nil {
		response.Error(c, http.StatusBadRequest, "LOGOUT_FAILED", err.Error())
		return
	}
	response.NoContent(c)
}
