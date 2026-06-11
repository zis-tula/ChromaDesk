package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error codes returned in API responses as {"code": "...", "error": "..."}
// Clients use `code` for i18n display; `error` is for server-side logging only.
const (
	// Generic
	CodeInvalidRequest  = "INVALID_REQUEST"
	CodeInternalError   = "INTERNAL_ERROR"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeNotFound        = "NOT_FOUND"

	// Auth
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeAccountDisabled    = "ACCOUNT_DISABLED"
	CodeSessionFailed      = "SESSION_FAILED"
	CodeTokenExpired       = "TOKEN_EXPIRED"

	// User
	CodeUsernameExists  = "USERNAME_EXISTS"
	CodePhoneExists     = "PHONE_EXISTS"
	CodePhoneNotFound   = "PHONE_NOT_FOUND"
	CodePhoneInvalid    = "PHONE_INVALID"
	CodePasswordWeak    = "PASSWORD_WEAK"
	CodePasswordWrong   = "PASSWORD_WRONG"
	CodeUserNotFound    = "USER_NOT_FOUND"

	// SMS
	CodeSmsDisabled     = "SMS_DISABLED"
	CodeSmsRateLimit    = "SMS_RATE_LIMIT"
	CodeSmsDailyLimit   = "SMS_DAILY_LIMIT"
	CodeSmsInvalid      = "SMS_CODE_INVALID"
	CodeSmsExpired      = "SMS_CODE_EXPIRED"
	CodeSmsMaxAttempts  = "SMS_MAX_ATTEMPTS"
	CodeSmsSceneInvalid = "SMS_SCENE_INVALID"
)

// apiError writes a JSON error response with both code and error message.
func apiError(c *gin.Context, status int, code, msg string) {
	c.JSON(status, gin.H{"code": code, "error": msg})
}

// apiErrorBadRequest is a shorthand for 400 errors.
func apiErrorBadRequest(c *gin.Context, code, msg string) {
	apiError(c, http.StatusBadRequest, code, msg)
}
