package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceDataStatus struct {
	IsSyncing           bool   `json:"isSyncing,omitempty"`
	IsConnected         bool   `json:"isConnected,omitempty"`
	IsStarting          bool   `json:"isStarting,omitempty"`
	IsStarted           bool   `json:"isStarted,omitempty"`
	IsConflicted        bool   `json:"isConflicted,omitempty"`
	IsLoading           bool   `json:"isLoading,omitempty"`
	IsOnChatPage        bool   `json:"isOnChatPage,omitempty"`
	EnteredQrCodePageAt string `json:"enteredQrCodePageAt,omitempty"`
	DisconnectedAt      string `json:"disconnectedAt,omitempty"`
	IsOnQrPage          bool   `json:"isOnQrPage,omitempty"`
	IsQrCodeExpired     bool   `json:"isQrCodeExpired,omitempty"`
	IsWebConnected      bool   `json:"isWebConnected,omitempty"`
	IsWebSyncing        bool   `json:"isWebSyncing,omitempty"`
	Mode                string `json:"mode,omitempty"`
	MyId                string `json:"myId,omitempty"`
	MyName              string `json:"myName,omitempty"`
	MyNumber            string `json:"myNumber,omitempty"`
	QrCodeExpiresAt     string `json:"qrCodeExpiresAt,omitempty"`
	QrCodeUrl           string `json:"qrCodeUrl,omitempty"`
	State               string `json:"state,omitempty"`
	WaVersion           string `json:"waVersion,omitempty"`
}

type ServiceData struct {
	MyId                 string            `json:"myId,omitempty"`
	LastSyncAt           *time.Time        `json:"lastSyncAt,omitempty"`
	LastMessageTimestamp *time.Time        `json:"lastMessageTimestamp,omitempty"`
	SyncCount            int               `json:"syncCount,omitempty"`
	SyncFlowDone         bool              `json:"syncFlowDone,omitempty"`
	Status               ServiceDataStatus `json:"status,omitempty"`
	StatusTimestamp      *time.Time        `json:"statusTimestamp,omitempty"`
}

type ServiceSettings struct {
	KeepOnline        bool   `json:"keepOnline,omitempty"`
	WppConnectVersion string `json:"wppConnectVersion,omitempty"`
	WaVersion         string `json:"waVersion,omitempty"`
}

type Service struct {
	Id          uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name        string           `gorm:"type:string;not null"`
	Data        *ServiceData     `gorm:"type:jsonb;not null;default:'{}';serializer:json"`
	Settings    *ServiceSettings `gorm:"type:jsonb;not null;default:'{}';serializer:json"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	AccountId   uuid.UUID
	ServerPodId *uuid.UUID
	Account     *Account   `gorm:"foreignKey:AccountId"`
	ServerPod   *ServerPod `gorm:"foreignKey:ServerPodId"`
}

func (Service) TableName() string {
	return "services"
}
