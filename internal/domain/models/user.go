package models

import "time"

type User struct {
	ID             int
	Username       string
	Email          string
	HashedPass     []byte
	Avatar         string
	Role           string
	Last_login     time.Time
	Login_attempts int
	Account_locked bool
	Created_at     time.Time
}

type IsAdmin struct {
	IsAdmin bool
}
