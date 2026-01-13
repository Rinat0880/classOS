package classosbackend

import "time"

type DeviceStatus struct {
	DeviceName    string    `json:"device_name" db:"device_name"`
	Username      string    `json:"username" db:"username"`
	LastHeartbeat time.Time `json:"last_heartbeat" db:"last_heartbeat"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	IsOnline      bool      `json:"is_online" db:"-"`
}

type UserLog struct {
	ID         int       `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	DeviceName string    `json:"device_name" db:"device_name"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	LogType    string    `json:"log_type" db:"log_type"`
	Program    string    `json:"program" db:"program"`
	Action     string    `json:"action" db:"action"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type LogsFilter struct {
	Username   string
	DeviceName string
	FromDate   *time.Time
	ToDate     *time.Time
	Limit      int
	Offset     int
}
