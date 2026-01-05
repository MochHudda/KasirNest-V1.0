package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	UserID    string    `json:"user_id" firestore:"user_id"`
	Email     string    `json:"email" firestore:"email"`
	Name      string    `json:"name" firestore:"name"`
	Role      string    `json:"role" firestore:"role"` // "admin" or "kasir"
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	LastLogin time.Time `json:"last_login" firestore:"last_login"`
}

// UserRole constants
const (
	RoleAdmin = "admin"
	RoleKasir = "kasir"
)

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsKasir checks if user has kasir role
func (u *User) IsKasir() bool {
	return u.Role == RoleKasir
}
