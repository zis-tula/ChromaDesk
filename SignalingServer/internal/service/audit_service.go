package service

import (
	"context"
	"time"

	"chromadesk/signaling/internal/models"

	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

func (s *AuditService) Log(ctx context.Context, adminID uint, adminUsername, action, resourceType, resourceID, details, ip string) {
	entry := &models.AuditLog{
		AdminID:       adminID,
		AdminUsername: adminUsername,
		Action:        action,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		Details:       details,
		IP:            ip,
		CreatedAt:     time.Now(),
	}
	s.db.WithContext(ctx).Create(entry)
}

func (s *AuditService) List(ctx context.Context, offset, limit int, sort, order, action, adminUsername, dateFrom, dateTo string) ([]models.AuditLog, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.AuditLog{})

	if action != "" {
		query = query.Where("action = ?", action)
	}
	if adminUsername != "" {
		query = query.Where("admin_username LIKE ?", "%"+adminUsername+"%")
	}
	if dateFrom != "" {
		if t, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if dateTo != "" {
		if t, err := time.Parse(time.RFC3339, dateTo); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]bool{"created_at": true, "action": true, "admin_username": true}
	if !allowedSorts[sort] {
		sort = "created_at"
	}
	if order != "asc" {
		order = "desc"
	}

	var logs []models.AuditLog
	err := query.Order(sort + " " + order).Offset(offset).Limit(limit).Find(&logs).Error
	return logs, total, err
}
