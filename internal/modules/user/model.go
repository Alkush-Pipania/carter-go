package user

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        string
	Email     string
	Username  string
	Image     string
	Password  string
	Verified  bool
	CreatedAt time.Time
}

type InputCreateUser struct {
	Email        string
	Username     pgtype.Text
	PasswordHash string
}
