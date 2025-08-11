package util

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringArray []string

// Scan implements the sql.Scanner interface for reading from DB
func (s *StringArray) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, s)
}

// Value implements the driver.Valuer interface for saving to DB
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}
