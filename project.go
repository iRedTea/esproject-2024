package esproject

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringSlice []string
type Int64Slice []int64

type Project struct {
	Id      int64       `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name    string      `json:"name"`
	OwnerId int64       `json:"owner_id"`
	Tags    StringSlice `json:"tags" gorm:"type:text"`
	Members Int64Slice  `json:"members" gorm:"type:text"`
}

type ProjectInvite struct {
	Id      int64 `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Project int64 `json:"project"`
	User    int64 `json:"user"`
}

// Value Implementing the driver.Valuer interface for StringSlice
func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan Implementing the sql.Scanner interface for StringSlice
func (s *StringSlice) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), s); err != nil {
		return fmt.Errorf("failed to unmarshal StringSlice: %w", err)
	}
	return nil
}

// Value Implementing the driver.Valuer interface for Int64Slice
func (s Int64Slice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan Implementing the sql.Scanner interface for Int64Slice
func (s *Int64Slice) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), s); err != nil {
		return fmt.Errorf("failed to unmarshal Int64Slice: %w", err)
	}
	return nil
}
