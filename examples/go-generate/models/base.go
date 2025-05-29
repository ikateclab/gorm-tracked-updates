package models

import (
	"database/sql/driver"
	"errors"

	"github.com/bytedance/sonic"
)

// JsonbStringSlice represents a jsonb type in PostgreSQL
type JsonbStringSlice []string

// Scan implements the Scanner interface.
func (j *JsonbStringSlice) Scan(value interface{}) error {
	if value == nil {
		*j = JsonbStringSlice{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return sonic.Unmarshal(bytes, j)
}

// Value implements the driver Valuer interface.
func (j JsonbStringSlice) Value() (driver.Value, error) {
	return sonic.Marshal(j)
}
