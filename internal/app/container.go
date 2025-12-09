package app

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/jackc/pgx/v5"
)

type Container struct {
	DB *pgx.Conn
}

func NewContainer(ctx context.Context, db *db.Queries) *Container {

}
