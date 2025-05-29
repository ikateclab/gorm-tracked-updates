package models

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type AccountSettings struct {
}

type AccountData struct {
}

type Account struct {
	Id            uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name          string           `gorm:"type:string;not null"`
	Settings      *AccountSettings `gorm:"type:jsonb;not null;default:'{}';serializer:json"`
	Data          *AccountData     `gorm:"type:jsonb;not null;default:'{}';serializer:json"`
	IsActive      bool             `gorm:"default:false"`
	CorrelationId string           `gorm:"type:string;default:''"`
	WebhookUrl    string           `gorm:"type:text;not null;default:''"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Services      []*Service     `gorm:"foreignKey:AccountId"`
}

func (Account) TableName() string {
	return "accounts"
}
