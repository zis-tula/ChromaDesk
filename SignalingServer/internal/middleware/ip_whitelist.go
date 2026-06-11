package middleware

import (
	"net"
	"net/http"
	"strings"

	"chromadesk/signaling/internal/service"

	"github.com/gin-gonic/gin"
)

func IPWhitelistMiddleware(settingsService *service.SettingsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		settings := settingsService.Get()
		whitelist := strings.TrimSpace(settings.AdminIPWhitelist)
		if whitelist == "" {
			c.Next()
			return
		}

		clientIP := net.ParseIP(c.ClientIP())
		if clientIP == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid client IP"})
			return
		}

		entries := strings.Split(whitelist, "\n")
		for _, entry := range entries {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				continue
			}

			if strings.Contains(entry, "/") {
				_, cidr, err := net.ParseCIDR(entry)
				if err == nil && cidr.Contains(clientIP) {
					c.Next()
					return
				}
			} else {
				if net.ParseIP(entry) != nil && entry == clientIP.String() {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "IP not allowed"})
	}
}
