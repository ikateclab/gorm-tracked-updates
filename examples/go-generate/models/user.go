package models

//go:generate go run ../../../cmd/gorm-gen

import "time"

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Age       int
	Profile   Profile   `gorm:"embedded"`
	Addresses []Address `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Profile represents user profile information
type Profile struct {
	Bio       string
	Avatar    string
	Verified  bool
	Settings  map[string]interface{} `gorm:"serializer:json"`
	Metadata  map[string]string      `gorm:"serializer:json"`
}
