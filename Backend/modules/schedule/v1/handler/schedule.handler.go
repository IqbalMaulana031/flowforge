package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"flowforge-api/middleware"
	"flowforge-api/modules/schedule/v1/service"
	"flowforge-api/resource"
	"flowforge-api/response"
)

type ScheduleHandler struct{ service service.ScheduleUseCase }

func NewScheduleHandler(service service.ScheduleUseCase) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}
func (h *ScheduleHandler) Create(c *gin.Context) {
	var req resource.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Create(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.GetString(middleware.UserIDKey), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "SCHEDULE_CREATE_FAILED", err.Error())
		return
	}
	response.Created(c, res)
}
func (h *ScheduleHandler) List(c *gin.Context) {
	res, err := h.service.List(c.Request.Context(), c.GetString(middleware.TenantIDKey))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "SCHEDULE_LIST_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *ScheduleHandler) Detail(c *gin.Context) {
	res, err := h.service.List(c.Request.Context(), c.GetString(middleware.TenantIDKey))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "SCHEDULE_DETAIL_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *ScheduleHandler) Update(c *gin.Context) {
	var req resource.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Update(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("id"), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "SCHEDULE_UPDATE_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *ScheduleHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, "SCHEDULE_DELETE_FAILED", err.Error())
		return
	}
	response.NoContent(c)
}
