package service

import (
	classosbackend "github.com/rinat0880/classOS_backend"
	"github.com/rinat0880/classOS_backend/pkg/repository"
)

type DeviceService struct {
	repo repository.Device
}

func NewDeviceService(repo repository.Device) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) UpsertDeviceStatus(device classosbackend.DeviceStatus) error {
	return s.repo.UpsertDeviceStatus(device)
}

func (s *DeviceService) GetAllDevices() ([]classosbackend.DeviceStatus, error) {
	return s.repo.GetAllDevices()
}

func (s *DeviceService) GetOnlineDevices() ([]classosbackend.DeviceStatus, error) {
	return s.repo.GetOnlineDevices()
}

func (s *DeviceService) GetDeviceByName(deviceName string) (classosbackend.DeviceStatus, error) {
	return s.repo.GetDeviceByName(deviceName)
}

func (s *DeviceService) DeleteDevice(deviceName string) error {
	return s.repo.DeleteDevice(deviceName)
}
