package db

type Repository struct {
	database *Queries
}

func NewPostgresRepository(store *Queries) *Repository {
	return &Repository{
		database: store,
	}
}
