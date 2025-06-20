package models

import (
	"github.com/google/uuid"
)

// Tag represents a simple tag for testing array of objects
// @jsonb
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Item represents an item for testing array of objects  
// @jsonb
type Item struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price,omitempty"`
}

// SimpleModel demonstrates JSONB array fields issue
type SimpleModel struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name  string    `gorm:"type:string;not null" json:"name"`
	Tags  []*Tag    `gorm:"type:jsonb;serializer:json;default:'[]'" json:"tags,omitempty"`
	Items []*Item   `gorm:"type:jsonb;serializer:json;default:'[]'" json:"items,omitempty"`
}

func (SimpleModel) TableName() string {
	return "simple_models"
}
