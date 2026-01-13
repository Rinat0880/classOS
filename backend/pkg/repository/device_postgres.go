package repository

import (
	"time"

	classosbackend "github.com/rinat0880/classOS_backend"
	"github.com/jmoiron/sqlx"
)

type DevicePostgres struct {
	db *sqlx.DB
}

func NewDevicePostgres(db *sqlx.DB) *DevicePostgres {
	return &DevicePostgres{db: db}
}

func (r *DevicePostgres) UpsertDeviceStatus(device classosbackend.DeviceStatus) error {
	query := `
		INSERT INTO device_status (device_name, username, last_heartbeat, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (device_name) 
		DO UPDATE SET 
			username = EXCLUDED.username,
			last_heartbeat = EXCLUDED.last_heartbeat,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(query, device.DeviceName, device.Username, device.LastHeartbeat, time.Now())
	return err
}

func (r *DevicePostgres) GetAllDevices() ([]classosbackend.DeviceStatus, error) {
	var devices []classosbackend.DeviceStatus
	query := `SELECT device_name, username, last_heartbeat, created_at, updated_at FROM device_status ORDER BY last_heartbeat DESC`
	
	err := r.db.Select(&devices, query)
	if err != nil {
		return nil, err
	}

	for i := range devices {
		devices[i].IsOnline = time.Since(devices[i].LastHeartbeat) < 2*time.Minute
	}

	return devices, nil
}

func (r *DevicePostgres) GetOnlineDevices() ([]classosbackend.DeviceStatus, error) {
	var devices []classosbackend.DeviceStatus
	query := `
		SELECT device_name, username, last_heartbeat, created_at, updated_at 
		FROM device_status 
		WHERE last_heartbeat > $1
		ORDER BY last_heartbeat DESC
	`
	
	threshold := time.Now().Add(-2 * time.Minute)
	err := r.db.Select(&devices, query, threshold)
	if err != nil {
		return nil, err
	}

	for i := range devices {
		devices[i].IsOnline = true
	}

	return devices, nil
}

func (r *DevicePostgres) GetDeviceByName(deviceName string) (classosbackend.DeviceStatus, error) {
	var device classosbackend.DeviceStatus
	query := `SELECT device_name, username, last_heartbeat, created_at, updated_at FROM device_status WHERE device_name = $1`
	
	err := r.db.Get(&device, query, deviceName)
	if err != nil {
		return device, err
	}

	device.IsOnline = time.Since(device.LastHeartbeat) < 2*time.Minute
	return device, nil
}

func (r *DevicePostgres) DeleteDevice(deviceName string) error {
	query := `DELETE FROM device_status WHERE device_name = $1`
	_, err := r.db.Exec(query, deviceName)
	return err
}
