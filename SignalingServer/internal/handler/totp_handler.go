package handler

import (
	"net/http"

	"chromadesk/signaling/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type TOTPHandler struct {
	adminService *service.AdminUserService
	db           *gorm.DB
}

func NewTOTPHandler(adminService *service.AdminUserService, db *gorm.DB) *TOTPHandler {
	return &TOTPHandler{adminService: adminService, db: db}
}

func (h *TOTPHandler) Setup2FA(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	id, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.adminService.GetAdminUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ChromaDesk",
		AccountName: user.Username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate TOTP"})
		return
	}

	h.db.Model(user).Update("totp_secret", key.Secret())

	c.JSON(http.StatusOK, gin.H{
		"secret": key.Secret(),
		"qr_uri": key.URL(),
	})
}

func (h *TOTPHandler) Verify2FA(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	id, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code required"})
		return
	}

	user, err := h.adminService.GetAdminUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.TOTPSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA not setup"})
		return
	}

	if !totp.Validate(req.Code, user.TOTPSecret) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	h.db.Model(user).Update("totp_enabled", true)
	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled"})
}

func (h *TOTPHandler) Disable2FA(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	id, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code required"})
		return
	}

	user, err := h.adminService.GetAdminUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if !user.TOTPEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA not enabled"})
		return
	}

	if !totp.Validate(req.Code, user.TOTPSecret) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	h.db.Model(user).Updates(map[string]interface{}{"totp_enabled": false, "totp_secret": ""})
	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled"})
}
