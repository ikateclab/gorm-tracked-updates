package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ServerPod struct {
	Id              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name            string         `gorm:"type:varchar(255);unique;not null"`
	Address         string         `gorm:"type:varchar(255);unique;not null"`
	Version         string         `gorm:"type:varchar(255);not null"`
	Settings        datatypes.JSON `gorm:"type:jsonb;default:'{}';not null"`
	LastPingAt      time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	ServerPodTypeId uuid.UUID
	ServerPodType   ServerPodType `gorm:"foreignKey:ServerPodTypeId"`
}

func (ServerPod) TableName() string {
	return "server_pods"
}
