package repository

import (
	"fmt"
	"strings"

	classosbackend "github.com/rinat0880/classOS_backend"
	"github.com/jmoiron/sqlx"
)

type LogsPostgres struct {
	db *sqlx.DB
}

func NewLogsPostgres(db *sqlx.DB) *LogsPostgres {
	return &LogsPostgres{db: db}
}

func (r *LogsPostgres) SaveLogs(logs []classosbackend.UserLog) error {
	if len(logs) == 0 {
		return nil
	}

	query := `
		INSERT INTO user_logs (username, device_name, timestamp, log_type, program, action)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, log := range logs {
		_, err = tx.Exec(query, log.Username, log.DeviceName, log.Timestamp, log.LogType, log.Program, log.Action)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *LogsPostgres) GetLogsByUsername(username string, limit, offset int) ([]classosbackend.UserLog, error) {
	var logs []classosbackend.UserLog
	query := `
		SELECT id, username, device_name, timestamp, log_type, program, action, created_at
		FROM user_logs
		WHERE username = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.Select(&logs, query, username, limit, offset)
	return logs, err
}

func (r *LogsPostgres) GetLogsByDevice(deviceName string, limit, offset int) ([]classosbackend.UserLog, error) {
	var logs []classosbackend.UserLog
	query := `
		SELECT id, username, device_name, timestamp, log_type, program, action, created_at
		FROM user_logs
		WHERE device_name = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.Select(&logs, query, deviceName, limit, offset)
	return logs, err
}

func (r *LogsPostgres) GetLogsFiltered(filter classosbackend.LogsFilter) ([]classosbackend.UserLog, error) {
	var logs []classosbackend.UserLog
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Username != "" {
		conditions = append(conditions, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, filter.Username)
		argIndex++
	}

	if filter.DeviceName != "" {
		conditions = append(conditions, fmt.Sprintf("device_name = $%d", argIndex))
		args = append(args, filter.DeviceName)
		argIndex++
	}

	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("timestamp >= $%d", argIndex))
		args = append(args, filter.FromDate)
		argIndex++
	}

	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("timestamp <= $%d", argIndex))
		args = append(args, filter.ToDate)
		argIndex++
	}

	query := `
		SELECT id, username, device_name, timestamp, log_type, program, action, created_at
		FROM user_logs
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY timestamp DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	err := r.db.Select(&logs, query, args...)
	return logs, err
}

func (r *LogsPostgres) GetLogsCount(filter classosbackend.LogsFilter) (int, error) {
	var count int
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Username != "" {
		conditions = append(conditions, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, filter.Username)
		argIndex++
	}

	if filter.DeviceName != "" {
		conditions = append(conditions, fmt.Sprintf("device_name = $%d", argIndex))
		args = append(args, filter.DeviceName)
		argIndex++
	}

	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("timestamp >= $%d", argIndex))
		args = append(args, filter.FromDate)
		argIndex++
	}

	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("timestamp <= $%d", argIndex))
		args = append(args, filter.ToDate)
		argIndex++
	}

	query := "SELECT COUNT(*) FROM user_logs"

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Get(&count, query, args...)
	return count, err
}
