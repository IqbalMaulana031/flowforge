package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"flowforge-api/middleware"
	"flowforge-api/modules/workflow/v1/service"
	"flowforge-api/resource"
	"flowforge-api/response"
)

type WorkflowHandler struct{ service service.WorkflowUseCase }

func NewWorkflowHandler(service service.WorkflowUseCase) *WorkflowHandler {
	return &WorkflowHandler{service: service}
}
func (h *WorkflowHandler) Create(c *gin.Context) {
	var req resource.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Create(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.GetString(middleware.UserIDKey), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "WORKFLOW_CREATE_FAILED", err.Error())
		return
	}
	response.Created(c, res)
}
func (h *WorkflowHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, meta, err := h.service.List(c.Request.Context(), c.GetString(middleware.TenantIDKey), page, limit, c.Query("q"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "WORKFLOW_LIST_FAILED", err.Error())
		return
	}
	response.Success(c, res, meta)
}
func (h *WorkflowHandler) Detail(c *gin.Context) {
	res, err := h.service.Detail(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *WorkflowHandler) Update(c *gin.Context) {
	var req resource.UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Update(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.GetString(middleware.UserIDKey), c.Param("id"), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "WORKFLOW_UPDATE_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *WorkflowHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, "WORKFLOW_DELETE_FAILED", err.Error())
		return
	}
	response.NoContent(c)
}
func (h *WorkflowHandler) Versions(c *gin.Context) {
	res, err := h.service.Versions(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "WORKFLOW_VERSIONS_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *WorkflowHandler) Rollback(c *gin.Context) {
	var req resource.RollbackWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.Rollback(c.Request.Context(), c.GetString(middleware.TenantIDKey), c.Param("id"), req.VersionNumber)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "WORKFLOW_ROLLBACK_FAILED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
