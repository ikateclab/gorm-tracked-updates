package models

// Address represents a user address
type Address struct {
	ID      uint   `gorm:"primaryKey"`
	UserID  uint   `gorm:"not null"`
	Type    string `gorm:"not null"` // home, work, etc.
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
	Primary bool `gorm:"default:false"`
}
