package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	classosbackend "github.com/rinat0880/classOS_backend"
)

type Authorization interface {
	CreateUser(user classosbackend.User) (int, error)
	GetUser(username, password string) (classosbackend.User, error)
}

type Group interface {
	Create(checkerId int, group classosbackend.Group) (int, error)
	GetAll(checkerId int) ([]classosbackend.Group, error)
	GetById(checkerId, groupId int) (classosbackend.Group, error)
	Delete(checkerId, groupId int) error
	Update(checkerId, groupId int, input classosbackend.UpdateGroupInput) error
	
	// Методы для транзакций
	BeginTransaction() (*sql.Tx, error)
	CreateWithTx(tx *sql.Tx, checkerId int, group classosbackend.Group) (int, error)
	UpdateWithTx(tx *sql.Tx, checkerId, groupId int, input classosbackend.UpdateGroupInput) error
	DeleteWithTx(tx *sql.Tx, checkerId, groupId int) error
}

type User interface {
	Create(groupId int, user classosbackend.User) (int, error)
	GetAll(checkerId int) ([]classosbackend.User, error)
	GetById(checkerId, userId int) (classosbackend.User, error)
	Delete(checkerId, userId int) error
	Update(checkerId, userId int, input classosbackend.UpdateUserInput) error
	
	// Методы для транзакций
	BeginTransaction() (*sql.Tx, error)
	CreateWithTx(tx *sql.Tx, groupId int, user classosbackend.User) (int, error)
	UpdateWithTx(tx *sql.Tx, checkerId, userId int, input classosbackend.UpdateUserInput) error
	DeleteWithTx(tx *sql.Tx, checkerId, userId int) error
}

type Device interface {
	UpsertDeviceStatus(device classosbackend.DeviceStatus) error
	GetAllDevices() ([]classosbackend.DeviceStatus, error)
	GetOnlineDevices() ([]classosbackend.DeviceStatus, error)
	GetDeviceByName(deviceName string) (classosbackend.DeviceStatus, error)
	DeleteDevice(deviceName string) error
}

type Logs interface {
	SaveLogs(logs []classosbackend.UserLog) error
	GetLogsByUsername(username string, limit, offset int) ([]classosbackend.UserLog, error)
	GetLogsByDevice(deviceName string, limit, offset int) ([]classosbackend.UserLog, error)
	GetLogsFiltered(filter classosbackend.LogsFilter) ([]classosbackend.UserLog, error)
	GetLogsCount(filter classosbackend.LogsFilter) (int, error)
}

type Repository struct {
	Authorization
	Group
	User
	Device
	Logs
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Group:         NewGroupPostgres(db),
		User:          NewUserPostgres(db),
		Device:        NewDevicePostgres(db),
		Logs:          NewLogsPostgres(db),
	}
}