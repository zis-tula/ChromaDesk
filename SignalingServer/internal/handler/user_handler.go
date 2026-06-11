package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"chromadesk/signaling/internal/models"
	"chromadesk/signaling/internal/service"
	"strconv"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler handles admin-facing user management CRUD operations.
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// validatePassword checks minimum password strength: ≥8 chars, has letter and digit.
func validatePassword(p string) error {
	if len(p) < 8 {
		return errors.New("密码至少8位")
	}
	hasLetter, hasDigit := false, false
	for _, r := range p {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("密码必须包含字母和数字")
	}
	return nil
}

// GetUsers handles GET /admin/user-list
// Supports: ?page=1&size=20&sort=created_at&order=desc&search=&level=&status=&channelType=
func (h *UserHandler) GetUsers(c *gin.Context) {
	p := ParsePagination(c)

	// Validate sort field
	allowedSorts := map[string]bool{
		"created_at": true, "updated_at": true, "username": true,
		"level": true, "device_count": true, "id": true,
	}
	if !allowedSorts[p.Sort] {
		p.Sort = "created_at"
	}

	query := h.db.Model(&models.User{})

	if p.Search != "" {
		like := "%" + p.Search + "%"
		query = query.Where("username LIKE ? OR phone LIKE ? OR email LIKE ?", like, like, like)
	}
	if level := c.Query("level"); level != "" {
		query = query.Where("level = ?", level)
	}
	if statusStr := c.Query("status"); statusStr == "true" || statusStr == "false" {
		query = query.Where("status = ?", statusStr == "true")
	}
	if channelType := c.Query("channelType"); channelType != "" {
		query = query.Where("channel_type = ?", channelType)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	query.Order(p.OrderClause()).Offset(p.Offset()).Limit(p.Size).Find(&users)

	type UserWithDevices struct {
		models.User
		Devices []models.UserDevice `json:"devices"`
	}

	result := make([]UserWithDevices, 0, len(users))
	for _, user := range users {
		var devices []models.UserDevice
		h.db.Where("user_id = ? AND status = ?", user.ID, true).Find(&devices)
		result = append(result, UserWithDevices{User: user, Devices: devices})
	}

	c.JSON(http.StatusOK, NewPaginatedResponse(result, total, p))
}

// GetUserDetail handles GET /admin/user-list/:id/details
func (h *UserHandler) GetUserDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if result := h.db.First(&user, id); result.Error != nil {
		apiError(c, http.StatusNotFound, CodeUserNotFound, "用户不存在")
		return
	}

	// Bound devices
	var devices []models.UserDevice
	h.db.Where("user_id = ? AND status = ?", user.ID, true).Find(&devices)

	// Connection history
	var histories []models.ConnectionHistory
	h.db.Where("user_id = ?", user.ID).Order("created_at DESC").Limit(50).Find(&histories)

	historyResult := make([]gin.H, 0, len(histories))
	for _, hist := range histories {
		historyResult = append(historyResult, gin.H{
			"id":         hist.ID,
			"device_id":  hist.DeviceID,
			"status":     hist.Status,
			"duration":   hist.Duration,
			"connect_ip": hist.ConnectIP,
			"created_at": hist.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user":              userJSON(&user),
		"devices":           devices,
		"connectionHistory": historyResult,
	})
}

// GetUser handles GET /admin/user-list/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if result := h.db.First(&user, id); result.Error != nil {
		apiError(c, http.StatusNotFound, CodeUserNotFound, "用户不存在")
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST /admin/user-list
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Password    string `json:"password" binding:"required"`
		Level       string `json:"level"`
		ChannelType string `json:"channelType"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, err.Error())
		return
	}

	var existing models.User
	if result := h.db.Where("username = ?", req.Username).First(&existing); result.Error == nil {
		apiErrorBadRequest(c, CodeUsernameExists, "用户名已存在")
		return
	}

	if err := validatePassword(req.Password); err != nil {
		apiErrorBadRequest(c, CodePasswordWeak, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		apiError(c, http.StatusInternalServerError, CodeInternalError, "密码加密失败")
		return
	}

	level := req.Level
	if level == "" {
		level = "V1"
	}
	channelType := req.ChannelType
	if channelType == "" {
		channelType = "全球"
	}

	user := models.User{
		Username:    req.Username,
		Phone:       req.Phone,
		Email:       req.Email,
		Password:    string(hashedPassword),
		Level:       level,
		ChannelType: channelType,
		Status:      true,
	}

	if result := h.db.Create(&user); result.Error != nil {
		apiError(c, http.StatusInternalServerError, CodeInternalError, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户创建成功", "user": user})
}

// UpdateUser handles PUT /admin/user-list/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if result := h.db.First(&user, id); result.Error != nil {
		apiError(c, http.StatusNotFound, CodeUserNotFound, "用户不存在")
		return
	}

	var req struct {
		Username    string `json:"username"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Level       string `json:"level"`
		DeviceCount int    `json:"deviceCount"`
		ChannelType string `json:"channelType"`
		Status      *bool  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, err.Error())
		return
	}

	if req.Username != "" {
		var existing models.User
		if h.db.Where("username = ? AND id != ?", req.Username, id).First(&existing).Error == nil {
			apiErrorBadRequest(c, CodeUsernameExists, "用户名已存在")
			return
		}
		user.Username = req.Username
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		if err := validatePassword(req.Password); err != nil {
			apiErrorBadRequest(c, CodePasswordWeak, err.Error())
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			apiError(c, http.StatusInternalServerError, CodeInternalError, "密码加密失败")
			return
		}
		user.Password = string(hashed)
	}
	if req.Level != "" {
		user.Level = req.Level
	}
	if req.ChannelType != "" {
		user.ChannelType = req.ChannelType
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	user.DeviceCount = req.DeviceCount

	if result := h.db.Save(&user); result.Error != nil {
		apiError(c, http.StatusInternalServerError, CodeInternalError, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户更新成功", "user": user})
}

// DeleteUser handles DELETE /admin/user-list/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	// Delete all rows referencing this user before deleting the user
	for _, model := range []interface{}{&models.UserFavorite{}, &models.UserDevice{}, &models.ConnectionHistory{}} {
		if err := h.db.Where("user_id = ?", id).Delete(model).Error; err != nil {
			apiError(c, http.StatusInternalServerError, CodeInternalError, err.Error())
			return
		}
	}
	// Unbind devices (nullable FK — set to NULL rather than delete)
	if err := h.db.Model(&models.Device{}).Where("user_id = ?", id).Update("user_id", nil).Error; err != nil {
		apiError(c, http.StatusInternalServerError, CodeInternalError, err.Error())
		return
	}
	if result := h.db.Delete(&models.User{}, id); result.Error != nil {
		apiError(c, http.StatusInternalServerError, CodeInternalError, result.Error.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// UpdateUserDeviceCount handles PUT /admin/user-list/:id/device-count
func (h *UserHandler) UpdateUserDeviceCount(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		DeviceCount int `json:"deviceCount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, err.Error())
		return
	}

	if result := h.db.Model(&models.User{}).Where("id = ?", id).Update("device_count", req.DeviceCount); result.Error != nil {
		apiError(c, http.StatusInternalServerError, CodeInternalError, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备数量更新成功"})
}

// ---------------------------------------------------------------------------
// UserAuth: Redis-backed token authentication for end-users (7-day TTL).
// ---------------------------------------------------------------------------

const userTokenTTL = 7 * 24 * time.Hour

// UserAuth manages user session tokens in Redis.
type UserAuth struct {
	db             *gorm.DB
	rdb            *redis.Client
	sms            *service.SmsService
	logoutNotifier func(userID uint, deviceID string)
}

// NewUserAuth creates a new UserAuth instance.
func NewUserAuth(db *gorm.DB, rdb *redis.Client) *UserAuth {
	return &UserAuth{db: db, rdb: rdb}
}

// SetLogoutNotifier sets the callback to notify other devices when a device logs out.
func (a *UserAuth) SetLogoutNotifier(fn func(userID uint, deviceID string)) {
	a.logoutNotifier = fn
}

// SetSmsService injects the SMS service (may be nil if SMS is disabled).
func (a *UserAuth) SetSmsService(sms *service.SmsService) {
	a.sms = sms
}

// CleanupLoop is a no-op kept for API compatibility. Redis TTL handles expiry automatically.
func (a *UserAuth) CleanupLoop() {}

func (a *UserAuth) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (a *UserAuth) redisKey(token string) string {
	return fmt.Sprintf("user_token:%s", token)
}

// Register handles POST /api/v1/user/register
func (a *UserAuth) Register(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Phone       string `json:"phone"`
		SmsCode     string `json:"sms_code"`
		Email       string `json:"email"`
		Level       string `json:"level"`
		ChannelType string `json:"channelType"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "invalid request")
		return
	}

	if err := validatePassword(req.Password); err != nil {
		apiErrorBadRequest(c, CodePasswordWeak, err.Error())
		return
	}

	// When SMS is enabled and phone is provided, verify SMS code
	smsEnabled := a.sms != nil && a.sms.IsEnabled()
	if smsEnabled && req.Phone != "" {
		if req.SmsCode == "" {
			apiErrorBadRequest(c, CodeInvalidRequest, "提供手机号时验证码为必填项")
			return
		}
		if !service.ValidatePhone(req.Phone) {
			apiErrorBadRequest(c, CodePhoneInvalid, "手机号格式不正确")
			return
		}
		if err := a.sms.VerifyCode(c.Request.Context(), req.Phone, req.SmsCode); err != nil {
			apiErrorBadRequest(c, CodeSmsInvalid, err.Error())
			return
		}
	}

	var existing models.User
	if a.db.Where("username = ?", req.Username).First(&existing).Error == nil {
		apiErrorBadRequest(c, CodeUsernameExists, "用户名已存在")
		return
	}

	if req.Phone != "" {
		var phoneUser models.User
		if a.db.Where("phone = ?", req.Phone).First(&phoneUser).Error == nil {
			apiError(c, http.StatusConflict, CodePhoneExists, "该手机号已注册")
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		apiError(c, http.StatusInternalServerError, CodeInternalError, "密码加密失败")
		return
	}

	user := models.User{
		Username:    req.Username,
		Phone:       req.Phone,
		Email:       req.Email,
		Password:    string(hashedPassword),
		Level:       "V1",
		ChannelType: "全球",
		Status:      true,
	}

	if result := a.db.Create(&user); result.Error != nil {
		apiError(c, http.StatusInternalServerError, CodeInternalError, "创建用户失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"user":    userJSON(&user),
	})
}

// Login handles POST /api/v1/user/login
func (a *UserAuth) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "invalid request")
		return
	}

	var user models.User
	if a.db.Where("username = ?", req.Username).First(&user).Error != nil {
		apiError(c, http.StatusUnauthorized, CodeInvalidCredentials, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		apiError(c, http.StatusUnauthorized, CodeInvalidCredentials, "用户名或密码错误")
		return
	}

	if !user.Status {
		apiError(c, http.StatusForbidden, CodeAccountDisabled, "账号已被禁用")
		return
	}

	token := a.generateToken()
	if err := a.rdb.Set(context.Background(), a.redisKey(token), user.ID, userTokenTTL).Err(); err != nil {
		apiError(c, http.StatusInternalServerError, CodeSessionFailed, "session 存储失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": userJSON(&user)})
}

// LoginWithSms handles POST /api/v1/user/login-sms
func (a *UserAuth) LoginWithSms(c *gin.Context) {
	if a.sms == nil || !a.sms.IsEnabled() {
		apiError(c, http.StatusServiceUnavailable, CodeSmsDisabled, "短信服务未启用")
		return
	}

	var req struct {
		Phone   string `json:"phone" binding:"required"`
		SmsCode string `json:"sms_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "请提供手机号和验证码")
		return
	}

	if !service.ValidatePhone(req.Phone) {
		apiErrorBadRequest(c, CodePhoneInvalid, "手机号格式不正确")
		return
	}

	if err := a.sms.VerifyCode(c.Request.Context(), req.Phone, req.SmsCode); err != nil {
		apiError(c, http.StatusUnauthorized, CodeSmsInvalid, err.Error())
		return
	}

	var user models.User
	if a.db.Where("phone = ?", req.Phone).First(&user).Error != nil {
		apiError(c, http.StatusUnauthorized, CodePhoneNotFound, "该手机号未注册")
		return
	}

	if !user.Status {
		apiError(c, http.StatusForbidden, CodeAccountDisabled, "账号已被禁用")
		return
	}

	token := a.generateToken()
	if err := a.rdb.Set(context.Background(), a.redisKey(token), user.ID, userTokenTTL).Err(); err != nil {
		apiError(c, http.StatusInternalServerError, CodeSessionFailed, "session 存储失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": userJSON(&user)})
}

// AuthRequired returns a Gin middleware that requires a valid user token.
func (a *UserAuth) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": CodeUnauthorized, "error": "未授权"})
			return
		}

		val, err := a.rdb.Get(context.Background(), a.redisKey(token)).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": CodeTokenExpired, "error": "token已过期"})
			return
		}

		userID, _ := strconv.ParseUint(val, 10, 64)
		c.Set("authed_user_id", uint(userID))
		c.Next()
	}
}

// GetUserIDFromToken extracts the user ID from the request token.
func (a *UserAuth) GetUserIDFromToken(c *gin.Context) uint {
	token := extractToken(c)
	if token == "" {
		return 0
	}
	val, err := a.rdb.Get(context.Background(), a.redisKey(token)).Result()
	if err != nil {
		return 0
	}
	userID, _ := strconv.ParseUint(val, 10, 64)
	return uint(userID)
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return c.Query("token")
}

// Logout handles POST /api/v1/user/logout
func (a *UserAuth) Logout(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		apiErrorBadRequest(c, CodeInvalidRequest, "token不能为空")
		return
	}

	// Get user ID from token before deleting it
	val, err := a.rdb.Get(context.Background(), a.redisKey(token)).Result()
	if err == nil {
		userID, _ := strconv.ParseUint(val, 10, 64)
		if userID > 0 {
			// Get device_id from request (optional query param or JSON body)
			deviceID := c.Query("device_id")
			if deviceID == "" {
				var body struct {
					DeviceID string `json:"device_id"`
				}
				c.ShouldBindJSON(&body)
				deviceID = body.DeviceID
			}
			if deviceID != "" {
				// Clear logged_in for this specific device
				a.db.Model(&models.Device{}).Where("device_id = ? AND user_id = ?", deviceID, userID).Update("logged_in", false)
				if a.logoutNotifier != nil {
					a.logoutNotifier(uint(userID), deviceID)
				}
			}
		}
	}

	a.rdb.Del(context.Background(), a.redisKey(token))
	c.JSON(http.StatusOK, gin.H{"message": "退出登录成功"})
}

// GetMe handles GET /api/v1/user/me
func (a *UserAuth) GetMe(c *gin.Context) {
	userID, _ := c.Get("authed_user_id")
	var user models.User
	if a.db.First(&user, userID).Error != nil {
		apiError(c, http.StatusNotFound, CodeUserNotFound, "用户不存在")
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userJSON(&user)})
}

// userJSON returns a consistent user map for API responses.
func userJSON(u *models.User) gin.H {
	return gin.H{
		"id":          u.ID,
		"username":    u.Username,
		"phone":       u.Phone,
		"email":       u.Email,
		"level":       u.Level,
		"deviceCount": u.DeviceCount,
		"channelType": u.ChannelType,
		"status":      u.Status,
		"createdAt":   u.CreatedAt,
		"updatedAt":   u.UpdatedAt,
	}
}

// ChangePassword handles PUT /api/v1/user/password
func (a *UserAuth) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("authed_user_id")
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "invalid request")
		return
	}
	if err := validatePassword(req.NewPassword); err != nil {
		apiErrorBadRequest(c, CodePasswordWeak, err.Error())
		return
	}
	var user models.User
	if a.db.First(&user, userID).Error != nil {
		apiError(c, http.StatusNotFound, CodeUserNotFound, "用户不存在")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)) != nil {
		apiError(c, http.StatusUnauthorized, CodePasswordWrong, "原密码错误")
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	a.db.Model(&user).Update("password", string(hashed))
	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// SendResetPasswordCode handles POST /api/v1/user/reset-password (public, no auth)
func (a *UserAuth) SendResetPasswordCode(c *gin.Context) {
	if a.sms == nil || !a.sms.IsEnabled() {
		apiError(c, http.StatusServiceUnavailable, CodeSmsDisabled, "短信服务未启用")
		return
	}
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "invalid request")
		return
	}
	if !service.ValidatePhone(req.Phone) {
		apiErrorBadRequest(c, CodePhoneInvalid, "手机号格式不正确")
		return
	}
	var user models.User
	if a.db.Where("phone = ?", req.Phone).First(&user).Error != nil {
		apiError(c, http.StatusNotFound, CodePhoneNotFound, "该手机号未注册")
		return
	}
	if err := a.sms.SendCode(c.Request.Context(), req.Phone); err != nil {
		errMsg := err.Error()
		code, status := CodeInternalError, http.StatusInternalServerError
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

// ResetPassword handles PUT /api/v1/user/reset-password (public, no auth)
func (a *UserAuth) ResetPassword(c *gin.Context) {
	if a.sms == nil || !a.sms.IsEnabled() {
		apiError(c, http.StatusServiceUnavailable, CodeSmsDisabled, "短信服务未启用")
		return
	}
	var req struct {
		Phone       string `json:"phone" binding:"required"`
		SmsCode     string `json:"sms_code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "invalid request")
		return
	}
	if !service.ValidatePhone(req.Phone) {
		apiErrorBadRequest(c, CodePhoneInvalid, "手机号格式不正确")
		return
	}
	if err := validatePassword(req.NewPassword); err != nil {
		apiErrorBadRequest(c, CodePasswordWeak, err.Error())
		return
	}
	if err := a.sms.VerifyCode(c.Request.Context(), req.Phone, req.SmsCode); err != nil {
		apiErrorBadRequest(c, CodeSmsInvalid, err.Error())
		return
	}
	var user models.User
	if a.db.Where("phone = ?", req.Phone).First(&user).Error != nil {
		apiError(c, http.StatusNotFound, CodePhoneNotFound, "该手机号未注册")
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	a.db.Model(&user).Update("password", string(hashed))
	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// ChangeUsername handles PUT /api/v1/user/username
func (a *UserAuth) ChangeUsername(c *gin.Context) {
	userID, _ := c.Get("authed_user_id")
	var req struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "invalid request")
		return
	}
	var existing models.User
	if a.db.Where("username = ? AND id != ?", req.Username, userID).First(&existing).Error == nil {
		apiErrorBadRequest(c, CodeUsernameExists, "用户名已存在")
		return
	}
	a.db.Model(&models.User{}).Where("id = ?", userID).Update("username", req.Username)
	c.JSON(http.StatusOK, gin.H{"message": "用户名修改成功"})
}

// ChangePhone handles PUT /api/v1/user/phone
func (a *UserAuth) ChangePhone(c *gin.Context) {
	if a.sms == nil || !a.sms.IsEnabled() {
		apiError(c, http.StatusServiceUnavailable, CodeSmsDisabled, "短信服务未启用")
		return
	}
	userID, _ := c.Get("authed_user_id")
	var req struct {
		Phone   string `json:"phone" binding:"required"`
		SmsCode string `json:"sms_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "invalid request")
		return
	}
	if !service.ValidatePhone(req.Phone) {
		apiErrorBadRequest(c, CodePhoneInvalid, "手机号格式不正确")
		return
	}
	var existing models.User
	if a.db.Where("phone = ? AND id != ?", req.Phone, userID).First(&existing).Error == nil {
		apiError(c, http.StatusConflict, CodePhoneExists, "该手机号已被使用")
		return
	}
	if err := a.sms.VerifyCode(c.Request.Context(), req.Phone, req.SmsCode); err != nil {
		apiErrorBadRequest(c, CodeSmsInvalid, err.Error())
		return
	}
	a.db.Model(&models.User{}).Where("id = ?", userID).Update("phone", req.Phone)
	c.JSON(http.StatusOK, gin.H{"message": "手机号修改成功"})
}

// ChangeEmail handles PUT /api/v1/user/email
func (a *UserAuth) ChangeEmail(c *gin.Context) {
	userID, _ := c.Get("authed_user_id")
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErrorBadRequest(c, CodeInvalidRequest, "invalid request")
		return
	}
	a.db.Model(&models.User{}).Where("id = ?", userID).Update("email", req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "邮箱修改成功"})
}
