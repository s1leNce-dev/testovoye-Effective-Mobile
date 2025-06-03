package repository

import "testovoye/internal/domain"

type Person interface {
	Create(newPerson domain.Person) error
	GetByFilter(query domain.PersonQuery) (domain.PaginatedPersons, error)
	GetByID(id string) (domain.Person, error)
	DeleteByID(id string) error
	Edit(editedPerson domain.Person) error
}
