package repository

import "testovoye/internal/domain"

type Person interface {
	Create(newPerson domain.Person) error
	Get(query domain.PersonQuery) (domain.PaginatedPersons, error)
	DeleteByID(id string) error
	Edit(editedPerson domain.Person) error
}
