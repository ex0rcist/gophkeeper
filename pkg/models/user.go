// Models used by server
package models

import "time"

// Beloved one
type User struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"-" db:"password"`
}
