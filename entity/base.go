package entity

import (
	"database/sql"
	"net"
	"time"

	"github.com/google/uuid"
)

type Timestamps struct {
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

type IPAddress struct {
	net.IP
}

type NullableUUID struct {
	UUID  uuid.UUID
	Valid bool
}

type NullableString struct {
	String string
	Valid  bool
}

type NullableTime struct {
	Time  time.Time
	Valid bool
}

type NullableInt struct {
	Int   int
	Valid bool
}
