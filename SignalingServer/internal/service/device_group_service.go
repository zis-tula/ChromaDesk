package service

import (
	"context"
	"errors"

	"chromadesk/signaling/internal/models"
	"chromadesk/signaling/internal/repository"
)

type DeviceGroupService struct {
	repo *repository.DeviceGroupRepository
}

func NewDeviceGroupService(repo *repository.DeviceGroupRepository) *DeviceGroupService {
	return &DeviceGroupService{repo: repo}
}

func (s *DeviceGroupService) Create(ctx context.Context, group *models.DeviceGroup) error {
	if group.Name == "" {
		return errors.New("name is required")
	}
	return s.repo.Create(ctx, group)
}

func (s *DeviceGroupService) GetAll(ctx context.Context) ([]models.DeviceGroup, error) {
	return s.repo.GetAll(ctx)
}

func (s *DeviceGroupService) GetByID(ctx context.Context, id uint) (*models.DeviceGroup, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DeviceGroupService) Update(ctx context.Context, group *models.DeviceGroup) error {
	return s.repo.Update(ctx, group)
}

func (s *DeviceGroupService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *DeviceGroupService) AddDevices(ctx context.Context, groupID uint, deviceIDs []string) error {
	return s.repo.AddDevices(ctx, groupID, deviceIDs)
}

func (s *DeviceGroupService) RemoveDevices(ctx context.Context, groupID uint, deviceIDs []string) error {
	return s.repo.RemoveDevices(ctx, groupID, deviceIDs)
}

func (s *DeviceGroupService) GetDeviceIDs(ctx context.Context, groupID uint) ([]string, error) {
	return s.repo.GetDeviceIDs(ctx, groupID)
}

func (s *DeviceGroupService) CountDevices(ctx context.Context, groupID uint) (int64, error) {
	return s.repo.CountDevices(ctx, groupID)
}
