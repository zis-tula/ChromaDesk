package repository

import (
	"context"
	"chromadesk/signaling/internal/models"
	"time"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// Create creates a new device
func (r *DeviceRepository) Create(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

// GetByDeviceID retrieves a device by device_id
func (r *DeviceRepository) GetByDeviceID(ctx context.Context, deviceID string) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&device).Error
	return &device, err
}

// GetByUUID retrieves a device by UUID
func (r *DeviceRepository) GetByUUID(ctx context.Context, uuid string) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).Where("device_uuid = ?", uuid).First(&device).Error
	return &device, err
}

// Update updates a device
func (r *DeviceRepository) Update(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

// SetOnline updates the online status of a device
func (r *DeviceRepository) SetOnline(ctx context.Context, deviceID string, online bool) error {
	updates := map[string]interface{}{"online": online}
	if online {
		updates["last_seen"] = time.Now()
	}
	return r.db.WithContext(ctx).Model(&models.Device{}).
		Where("device_id = ?", deviceID).
		Updates(updates).Error
}

// UpdateLastSeen updates the last_seen timestamp of a device
func (r *DeviceRepository) UpdateLastSeen(ctx context.Context, deviceID string) error {
	return r.db.WithContext(ctx).Model(&models.Device{}).
		Where("device_id = ?", deviceID).
		Update("last_seen", time.Now()).Error
}

// UpdateDeviceInfo updates OS, OSVersion, and AppVersion for a device
func (r *DeviceRepository) UpdateDeviceInfo(ctx context.Context, deviceID, os, osVersion, appVersion string) error {
	updates := map[string]interface{}{}
	if os != "" {
		updates["os"] = os
	}
	if osVersion != "" {
		updates["os_version"] = osVersion
	}
	if appVersion != "" {
		updates["app_version"] = appVersion
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&models.Device{}).
		Where("device_id = ?", deviceID).
		Updates(updates).Error
}

// List retrieves devices with pagination
func (r *DeviceRepository) List(ctx context.Context, offset, limit int) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&devices).Error
	return devices, err
}

// GetAll retrieves all devices
func (r *DeviceRepository) GetAll(ctx context.Context) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.WithContext(ctx).Find(&devices).Error
	return devices, err
}

// ListPaginated retrieves devices with pagination, search, filter, and sort
func (r *DeviceRepository) ListPaginated(ctx context.Context, offset, limit int, sort, order, search, os string, onlineFilter *bool) ([]models.Device, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Device{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("device_id LIKE ? OR device_name LIKE ?", like, like)
	}
	if os != "" {
		query = query.Where("os = ?", os)
	}
	if onlineFilter != nil {
		query = query.Where("online = ?", *onlineFilter)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var devices []models.Device
	orderClause := sort + " " + order
	err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&devices).Error
	return devices, total, err
}

// CountByDate counts devices created on or after the given date
func (r *DeviceRepository) CountSince(ctx context.Context, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Device{}).Where("created_at >= ?", since).Count(&count).Error
	return count, err
}
