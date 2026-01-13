package service

import (
	classosbackend "github.com/rinat0880/classOS_backend"
	"github.com/rinat0880/classOS_backend/pkg/repository"
)

type LogsService struct {
	repo repository.Logs
}

func NewLogsService(repo repository.Logs) *LogsService {
	return &LogsService{repo: repo}
}

func (s *LogsService) SaveLogs(logs []classosbackend.UserLog) error {
	return s.repo.SaveLogs(logs)
}

func (s *LogsService) GetLogsByUsername(username string, limit, offset int) ([]classosbackend.UserLog, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.repo.GetLogsByUsername(username, limit, offset)
}

func (s *LogsService) GetLogsByDevice(deviceName string, limit, offset int) ([]classosbackend.UserLog, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.repo.GetLogsByDevice(deviceName, limit, offset)
}

func (s *LogsService) GetLogsFiltered(filter classosbackend.LogsFilter) ([]classosbackend.UserLog, error) {
	if filter.Limit <= 0 {
		filter.Limit = 100
	}
	return s.repo.GetLogsFiltered(filter)
}

func (s *LogsService) GetLogsCount(filter classosbackend.LogsFilter) (int, error) {
	return s.repo.GetLogsCount(filter)
}
