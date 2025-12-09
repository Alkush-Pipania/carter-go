package user

import "time"

type User struct {
	ID        int
	Email     string
	Username  string
	Image     string
	Password  string
	SecretKey string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
