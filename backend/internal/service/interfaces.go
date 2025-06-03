package service

import "testovoye/internal/domain"

type Person interface {
	Add(newPerson NewPersonCreate) (string, error)
	GetWithFilter(query domain.PersonQuery) (domain.PaginatedPersons, error)
	GetByID(id string) (domain.Person, error)
	DeleteByID(id string) error
	Edit(editedPerson domain.Person) error
}
