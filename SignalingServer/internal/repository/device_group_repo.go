package repository

import (
	"context"

	"chromadesk/signaling/internal/models"

	"gorm.io/gorm"
)

type DeviceGroupRepository struct {
	db *gorm.DB
}

func NewDeviceGroupRepository(db *gorm.DB) *DeviceGroupRepository {
	return &DeviceGroupRepository{db: db}
}

func (r *DeviceGroupRepository) Create(ctx context.Context, group *models.DeviceGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *DeviceGroupRepository) GetByID(ctx context.Context, id uint) (*models.DeviceGroup, error) {
	var group models.DeviceGroup
	err := r.db.WithContext(ctx).First(&group, id).Error
	return &group, err
}

func (r *DeviceGroupRepository) GetAll(ctx context.Context) ([]models.DeviceGroup, error) {
	var groups []models.DeviceGroup
	err := r.db.WithContext(ctx).Order("name asc").Find(&groups).Error
	return groups, err
}

func (r *DeviceGroupRepository) Update(ctx context.Context, group *models.DeviceGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *DeviceGroupRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", id).Delete(&models.DeviceGroupMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.DeviceGroup{}, id).Error
	})
}

func (r *DeviceGroupRepository) AddDevices(ctx context.Context, groupID uint, deviceIDs []string) error {
	for _, did := range deviceIDs {
		member := models.DeviceGroupMember{GroupID: groupID, DeviceID: did}
		r.db.WithContext(ctx).Where("group_id = ? AND device_id = ?", groupID, did).FirstOrCreate(&member)
	}
	return nil
}

func (r *DeviceGroupRepository) RemoveDevices(ctx context.Context, groupID uint, deviceIDs []string) error {
	return r.db.WithContext(ctx).Where("group_id = ? AND device_id IN ?", groupID, deviceIDs).Delete(&models.DeviceGroupMember{}).Error
}

func (r *DeviceGroupRepository) GetDeviceIDs(ctx context.Context, groupID uint) ([]string, error) {
	var members []models.DeviceGroupMember
	err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&members).Error
	ids := make([]string, len(members))
	for i, m := range members {
		ids[i] = m.DeviceID
	}
	return ids, err
}

func (r *DeviceGroupRepository) GetGroupsByDevice(ctx context.Context, deviceID string) ([]models.DeviceGroup, error) {
	var groups []models.DeviceGroup
	err := r.db.WithContext(ctx).
		Joins("JOIN device_group_members ON device_group_members.group_id = device_groups.id").
		Where("device_group_members.device_id = ?", deviceID).
		Find(&groups).Error
	return groups, err
}

func (r *DeviceGroupRepository) CountDevices(ctx context.Context, groupID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.DeviceGroupMember{}).Where("group_id = ?", groupID).Count(&count).Error
	return count, err
}
