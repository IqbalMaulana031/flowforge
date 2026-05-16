package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"flowforge-api/middleware"
	"flowforge-api/modules/execution/v1/service"
	"flowforge-api/resource"
	"flowforge-api/response"
)

type RunHandler struct{ service service.RunUseCase }

func NewRunHandler(service service.RunUseCase) *RunHandler { return &RunHandler{service: service} }
func (h *RunHandler) Trigger(c *gin.Context) {
	var req resource.TriggerRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Trigger(c.Request.Context(), c.GetString(middleware.TenantIDKey), req.WorkflowID, req.Payload, "manual")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "TRIGGER_FAILED", err.Error())
		return
	}
	response.Created(c, res)
}
func (h *RunHandler) Webhook(c *gin.Context) {
	var payload map[string]any
	_ = c.ShouldBindJSON(&payload)
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "X-Tenant-ID header is required")
		return
	}
	res, err := h.service.Trigger(c.Request.Context(), tenantID, c.Param("workflowId"), payload, "webhook")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "WEBHOOK_FAILED", err.Error())
		return
	}
	response.Created(c, res)
}
func (h *RunHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, meta, err := h.service.List(c.Request.Context(), c.GetString(middleware.TenantIDKey), page, limit)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "RUN_LIST_FAILED", err.Error())
		return
	}
	response.Success(c, res, meta)
}
func (h *RunHandler) Detail(c *gin.Context) {
	res, err := h.service.Detail(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("runId"))
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *RunHandler) Steps(c *gin.Context) {
	res, err := h.service.Steps(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("runId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "RUN_STEPS_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *RunHandler) Logs(c *gin.Context) {
	res, err := h.service.Logs(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("runId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "RUN_LOGS_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *RunHandler) Cancel(c *gin.Context) {
	if err := h.service.Cancel(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("runId")); err != nil {
		response.Error(c, http.StatusBadRequest, "RUN_CANCEL_FAILED", err.Error())
		return
	}
	response.NoContent(c)
}
