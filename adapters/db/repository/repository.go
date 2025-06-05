package repository

import "github.com/medicue/adapters/db"

type Repository struct {
	database *db.Queries
}

func NewPostgresRepository(store *db.Queries) *Repository {
	return &Repository{
		database: store,
	}
}
