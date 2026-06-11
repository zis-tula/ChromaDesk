package handler

import (
	"fmt"
	"net/http"

	"chromadesk/signaling/internal/models"
	"chromadesk/signaling/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BatchHandler struct {
	db           *gorm.DB
	groupService *service.DeviceGroupService
}

func NewBatchHandler(db *gorm.DB, groupService *service.DeviceGroupService) *BatchHandler {
	return &BatchHandler{db: db, groupService: groupService}
}

func (h *BatchHandler) BatchDevices(c *gin.Context) {
	var req struct {
		Action  string   `json:"action" binding:"required"`
		IDs     []string `json:"ids" binding:"required"`
		GroupID uint     `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ids provided"})
		return
	}

	var affected int64
	var err error

	switch req.Action {
	case "delete":
		result := h.db.Where("device_id IN ?", req.IDs).Delete(&models.Device{})
		affected = result.RowsAffected
		err = result.Error
	case "group":
		if req.GroupID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "group_id required"})
			return
		}
		err = h.groupService.AddDevices(c.Request.Context(), req.GroupID, req.IDs)
		affected = int64(len(req.IDs))
	case "ungroup":
		if req.GroupID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "group_id required"})
			return
		}
		err = h.groupService.RemoveDevices(c.Request.Context(), req.GroupID, req.IDs)
		affected = int64(len(req.IDs))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d devices affected", affected)})
}

func (h *BatchHandler) BatchUsers(c *gin.Context) {
	var req struct {
		Action string `json:"action" binding:"required"`
		IDs    []uint `json:"ids" binding:"required"`
		Level  string `json:"level"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ids provided"})
		return
	}

	var result *gorm.DB

	switch req.Action {
	case "enable":
		result = h.db.Model(&models.User{}).Where("id IN ?", req.IDs).Update("status", true)
	case "disable":
		result = h.db.Model(&models.User{}).Where("id IN ?", req.IDs).Update("status", false)
	case "delete":
		result = h.db.Where("id IN ?", req.IDs).Delete(&models.User{})
	case "set-level":
		if req.Level == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "level required"})
			return
		}
		result = h.db.Model(&models.User{}).Where("id IN ?", req.IDs).Update("level", req.Level)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d users affected", result.RowsAffected)})
}
