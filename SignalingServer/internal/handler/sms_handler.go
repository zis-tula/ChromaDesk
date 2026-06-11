package handler

import (
	"net/http"
	"chromadesk/signaling/internal/models"
	"chromadesk/signaling/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SmsHandler handles SMS verification code endpoints.
type SmsHandler struct {
	sms *service.SmsService
	db  *gorm.DB
}

func NewSmsHandler(sms *service.SmsService, db *gorm.DB) *SmsHandler {
	return &SmsHandler{sms: sms, db: db}
}

// SendCode handles POST /api/v1/sms/send
func (h *SmsHandler) SendCode(c *gin.Context) {
	if !h.sms.IsEnabled() {
		apiError(c, http.StatusServiceUnavailable, CodeSmsDisabled, "短信服务未启用")
		return
	}

	var req struct {
		Phone string `json:"phone" binding:"required"`
		Scene string `json:"scene" binding:"required"` // "register", "login", "reset-password", "change-phone"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "请提供手机号和场景")
		return
	}

	if !service.ValidatePhone(req.Phone) {
		apiErrorBadRequest(c, CodePhoneInvalid, "手机号格式不正确")
		return
	}

	var existing models.User
	phoneExists := h.db.Where("phone = ?", req.Phone).First(&existing).Error == nil

	switch req.Scene {
	case "register":
		if phoneExists {
			apiError(c, http.StatusConflict, CodePhoneExists, "该手机号已注册")
			return
		}
	case "login", "reset-password":
		if !phoneExists {
			apiError(c, http.StatusConflict, CodePhoneNotFound, "该手机号未注册")
			return
		}
	case "change-phone":
		if phoneExists {
			apiError(c, http.StatusConflict, CodePhoneExists, "该手机号已被使用")
			return
		}
	default:
		apiErrorBadRequest(c, CodeSmsSceneInvalid, "无效的场景参数")
		return
	}

	if err := h.sms.SendCode(c.Request.Context(), req.Phone); err != nil {
		errMsg := err.Error()
		code := CodeInternalError
		status := http.StatusInternalServerError
		if errMsg == "发送太频繁，请稍后再试" {
			code, status = CodeSmsRateLimit, http.StatusTooManyRequests
		} else if errMsg == "今日验证码发送次数已达上限" {
			code, status = CodeSmsDailyLimit, http.StatusTooManyRequests
		}
		apiError(c, status, code, errMsg)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送", "expires_in": 300})
}
