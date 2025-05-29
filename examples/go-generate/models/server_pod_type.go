package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceVersion struct {
	WppConnectVersion string `json:"wppConnectVersion,omitempty"`
	WaVersion         string `json:"waVersion,omitempty"`
}

type ServerPodType struct {
	Id                 uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name               string           `gorm:"type:varchar(255);unique;not null"`
	Version            *ServiceVersion  `gorm:"type:jsonb;not null;default:'{}';serializer:json"`
	AutoScalable       bool             `gorm:"not null;default:false"`
	Cloud              string           `gorm:"type:varchar(255);not null"`
	ServerSize         string           `gorm:"type:varchar(255);not null"`
	MaxPerPod          int              `gorm:"not null;default:1"`
	Min                int              `gorm:"not null;default:0"`
	DesiredAvailable   int              `gorm:"not null;default:0"`
	StartPriority      int              `gorm:"not null;default:0"`
	AccountIdWhitelist JsonbStringSlice `gorm:"type:jsonb;default:'[]';not null"`
	ServiceIdWhitelist JsonbStringSlice `gorm:"type:jsonb;default:'[]';not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (ServerPodType) TableName() string {
	return "server_pod_types"
}
