package service

import (
	"testovoye/internal/repository"
	"testovoye/internal/utils"
)

type Service struct {
	Person Person
}

type Deps struct {
	Repo          *repository.Repositories
	FetcherExtAPI utils.FetcherAPIByName

	Domain string
}

func NewService(deps Deps) *Service {
	personService := NewPersonService(deps.Repo.Person, deps.FetcherExtAPI, deps.Domain)

	return &Service{
		Person: personService,
	}
}
