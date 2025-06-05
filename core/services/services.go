package services

import (
	"github.com/medicue/adapters/db/repository"
)

type ServicesHandler struct {
	repositoryService repository.Repository
}

func ServicesAdapter(repo repository.Repository) *ServicesHandler {
	return &ServicesHandler{
		repositoryService: repo,
	}
}
