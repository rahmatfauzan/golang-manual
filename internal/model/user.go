package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID            string       `db:"id" json:"id"`
	Username      string       `db:"username" json:"username"`
	Email         string       `db:"email" json:"email"`
	PasswordHash  string       `db:"password_hash" json:"-"`
	FullName      string       `db:"full_name" json:"full_name"`
	AvatarURL     *string      `db:"avatar_url" json:"avatar_url,omitempty"`
	Bio           *string      `db:"bio" json:"bio,omitempty"`
	IsActive      bool         `db:"is_active" json:"is_active"`
	EmailVerified bool         `db:"email_verified" json:"email_verified"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at" json:"-"`
}
