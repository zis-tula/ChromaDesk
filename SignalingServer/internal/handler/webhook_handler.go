package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chromadesk/signaling/internal/models"
	"chromadesk/signaling/internal/service"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
}

func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookService: webhookService}
}

func (h *WebhookHandler) GetWebhooks(c *gin.Context) {
	webhooks, err := h.webhookService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"webhooks": webhooks})
}

func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	var req struct {
		Name    string   `json:"name" binding:"required"`
		URL     string   `json:"url" binding:"required"`
		Secret  string   `json:"secret"`
		Events  []string `json:"events" binding:"required"`
		Enabled bool     `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventsJSON, _ := json.Marshal(req.Events)
	webhook := &models.Webhook{
		Name:    req.Name,
		URL:     req.URL,
		Secret:  req.Secret,
		Events:  string(eventsJSON),
		Enabled: req.Enabled,
	}

	if err := h.webhookService.Create(c.Request.Context(), webhook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	webhook, err := h.webhookService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "webhook not found"})
		return
	}

	var req struct {
		Name    string   `json:"name"`
		URL     string   `json:"url"`
		Secret  string   `json:"secret"`
		Events  []string `json:"events"`
		Enabled *bool    `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		webhook.Name = req.Name
	}
	if req.URL != "" {
		webhook.URL = req.URL
	}
	if req.Secret != "" {
		webhook.Secret = req.Secret
	}
	if req.Events != nil {
		eventsJSON, _ := json.Marshal(req.Events)
		webhook.Events = string(eventsJSON)
	}
	if req.Enabled != nil {
		webhook.Enabled = *req.Enabled
	}

	if err := h.webhookService.Update(c.Request.Context(), webhook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.webhookService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *WebhookHandler) TestWebhook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.webhookService.Test(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "test event sent"})
}
