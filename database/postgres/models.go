package dbpostgres

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Required
	Name string // A regular string field
	Age  uint8  // An unsigned 8-bit integer

	// Optional
	Email        *string        // A pointer to a string, allowing for null values
	Birthday     *time.Time     // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields

	// Example
	//nolint:unused
	ignored string // fields that aren't exported are ignored
}
