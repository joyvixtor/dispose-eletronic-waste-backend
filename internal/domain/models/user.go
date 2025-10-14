package models

import "time"

type User struct {
	Id        string     `json:"id" db:"id"`
	FullName  string     `json:"full_name" db:"full_name"`
	UFPEEmail string     `json:"ufpe_email" db:"ufpe_email"`
	Password  string     `json:"-" db:"password"`
	Workplace string     `json:"workplace" db:"workplace"`
	Role      string     `json:"role" db:"role"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin *time.Time `json:"last_login" db:"last_login"`
}
