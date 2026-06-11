package handler

import (
	"net/http"
	"strconv"

	"chromadesk/signaling/internal/models"
	"chromadesk/signaling/internal/service"

	"github.com/gin-gonic/gin"
)

type DeviceGroupHandler struct {
	groupService *service.DeviceGroupService
}

func NewDeviceGroupHandler(groupService *service.DeviceGroupService) *DeviceGroupHandler {
	return &DeviceGroupHandler{groupService: groupService}
}

func (h *DeviceGroupHandler) GetGroups(c *gin.Context) {
	groups, err := h.groupService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type groupWithCount struct {
		models.DeviceGroup
		DeviceCount int64 `json:"device_count"`
	}

	result := make([]groupWithCount, len(groups))
	for i, g := range groups {
		count, _ := h.groupService.CountDevices(c.Request.Context(), g.ID)
		result[i] = groupWithCount{DeviceGroup: g, DeviceCount: count}
	}

	c.JSON(http.StatusOK, gin.H{"groups": result})
}

func (h *DeviceGroupHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := &models.DeviceGroup{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}
	if err := h.groupService.Create(c.Request.Context(), group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *DeviceGroupHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	group, err := h.groupService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Description != "" {
		group.Description = req.Description
	}
	if req.Color != "" {
		group.Color = req.Color
	}

	if err := h.groupService.Update(c.Request.Context(), group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *DeviceGroupHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.groupService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *DeviceGroupHandler) AddDevices(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		DeviceIDs []string `json:"device_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.groupService.AddDevices(c.Request.Context(), uint(id), req.DeviceIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "devices added"})
}

func (h *DeviceGroupHandler) RemoveDevices(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		DeviceIDs []string `json:"device_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.groupService.RemoveDevices(c.Request.Context(), uint(id), req.DeviceIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "devices removed"})
}

func (h *DeviceGroupHandler) GetGroupDevices(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ids, err := h.groupService.GetDeviceIDs(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"device_ids": ids})
}
