package models

import "time"

// Order represents a customer order
type Order struct {
	ID          uint        `gorm:"primaryKey"`
	UserID      uint        `gorm:"not null"`
	User        *User       `gorm:"foreignKey:UserID"`
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	Total       float64
	Status      string
	ShippingAddress Address `gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddress  Address `gorm:"embedded;embeddedPrefix:billing_"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	Total     float64 `gorm:"not null"`
}
