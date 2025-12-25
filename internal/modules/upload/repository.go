package upload

import "github.com/Alkush-Pipania/carter-go/pkg/db"

type repository struct {
	dbQuerier *db.Queries
}

func NewRepo(db *db.Queries) *repository {
	return &repository{
		dbQuerier: db,
	}
}
