package handler

import (
	"net/http"

	"chromadesk/signaling/internal/service"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	p := ParsePagination(c)
	action := c.Query("action")
	admin := c.Query("admin")
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	logs, total, err := h.auditService.List(c.Request.Context(), p.Offset(), p.Size, p.Sort, p.Order, action, admin, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, NewPaginatedResponse(logs, total, p))
}
