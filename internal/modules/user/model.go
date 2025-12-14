package user

import "time"

type User struct {
	ID        string
	Email     string
	Username  string
	Image     string
	Password  string
	Verified  bool
	CreatedAt time.Time
}
