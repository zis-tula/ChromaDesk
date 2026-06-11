package handler

import (
	"net/http"
	"chromadesk/signaling/internal/service"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	service *service.SettingsService
}

func NewSettingsHandler(service *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

// GetPublicSettings returns non-sensitive settings for public access.
func (h *SettingsHandler) GetPublicSettings(c *gin.Context) {
	s := h.service.Get()
	c.JSON(http.StatusOK, gin.H{
		"siteEnabled": s.SiteEnabled,
		"siteName":    s.SiteName,
		"loginLogo":   s.LoginLogo,
		"smallLogo":   s.SmallLogo,
		"favicon":     s.Favicon,
	})
}

// GetSettings returns all settings including sensitive ones (admin only).
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	s := h.service.Get()
	c.JSON(http.StatusOK, s)
}

// settingsInput carries only the editable fields from the frontend form.
type settingsInput struct {
	SiteEnabled        *bool   `json:"siteEnabled"`
	SiteName           *string `json:"siteName"`
	LoginLogo          *string `json:"loginLogo"`
	SmallLogo          *string `json:"smallLogo"`
	Favicon            *string `json:"favicon"`
	TurnURLs           *string `json:"turnUrls"`
	TurnAuthSecret     *string `json:"turnAuthSecret"`
	TurnCredentialTTL  *int    `json:"turnCredentialTtl"`
	StunURLs           *string `json:"stunUrls"`
	APIKey             *string `json:"apiKey"`
	AllowedOrigins     *string `json:"allowedOrigins"`
	SmsAccessKeyID     *string `json:"smsAccessKeyId"`
	SmsAccessKeySecret *string `json:"smsAccessKeySecret"`
	SmsSignName        *string `json:"smsSignName"`
	SmsTemplateCode    *string `json:"smsTemplateCode"`
}

// UpdateSettings merges the JSON payload into the existing settings row, preserving
// fields the frontend didn't send (like ID / timestamps).
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var input settingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing := h.service.Get()

	if input.SiteEnabled != nil {
		existing.SiteEnabled = *input.SiteEnabled
	}
	if input.SiteName != nil {
		existing.SiteName = *input.SiteName
	}
	if input.LoginLogo != nil {
		existing.LoginLogo = *input.LoginLogo
	}
	if input.SmallLogo != nil {
		existing.SmallLogo = *input.SmallLogo
	}
	if input.Favicon != nil {
		existing.Favicon = *input.Favicon
	}
	if input.TurnURLs != nil {
		existing.TurnURLs = *input.TurnURLs
	}
	if input.TurnAuthSecret != nil {
		existing.TurnAuthSecret = *input.TurnAuthSecret
	}
	if input.TurnCredentialTTL != nil {
		existing.TurnCredentialTTL = *input.TurnCredentialTTL
	}
	if input.StunURLs != nil {
		existing.StunURLs = *input.StunURLs
	}
	if input.APIKey != nil {
		existing.APIKey = *input.APIKey
	}
	if input.AllowedOrigins != nil {
		existing.AllowedOrigins = *input.AllowedOrigins
	}
	if input.SmsAccessKeyID != nil {
		existing.SmsAccessKeyID = *input.SmsAccessKeyID
	}
	if input.SmsAccessKeySecret != nil {
		existing.SmsAccessKeySecret = *input.SmsAccessKeySecret
	}
	if input.SmsSignName != nil {
		existing.SmsSignName = *input.SmsSignName
	}
	if input.SmsTemplateCode != nil {
		existing.SmsTemplateCode = *input.SmsTemplateCode
	}

	if err := h.service.Save(&existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}
