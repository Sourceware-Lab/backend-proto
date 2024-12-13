package DBpostgres

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	// Required
	Name string // A regular string field
	Age  uint8  // An unsigned 8-bit integer

	// Optional
	Email        *string        // A pointer to a string, allowing for null values
	Birthday     *time.Time     // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields

	// Example
	ignored string // fields that aren't exported are ignored
}
