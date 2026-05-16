package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"flowforge-api/modules/ai/v1/service"
	"flowforge-api/resource"
	"flowforge-api/response"
)

type AIHandler struct{ service service.AIUseCase }

func NewAIHandler(service service.AIUseCase) *AIHandler { return &AIHandler{service: service} }
func (h *AIHandler) GenerateWorkflow(c *gin.Context) {
	var req resource.GenerateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := h.service.GenerateWorkflow(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, "AI_NOT_CONFIGURED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
func (h *AIHandler) AnalyzeFailure(c *gin.Context) {
	var req resource.AnalyzeFailureRequest
	_ = c.ShouldBindJSON(&req)
	res, err := h.service.AnalyzeFailure(c.Request.Context(), c.Param("runId"), req)
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, "AI_NOT_CONFIGURED", err.Error())
		return
	}
	response.Success(c, res, nil)
}
