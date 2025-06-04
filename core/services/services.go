package services

import (
	"github.com/medicue/adapters/db"
)

type ServicesHandler struct {
	repositoryService db.Repository
}

func ServicesAdapter(repo db.Repository) *ServicesHandler {
	return &ServicesHandler{
		repositoryService: repo,
	}
}
