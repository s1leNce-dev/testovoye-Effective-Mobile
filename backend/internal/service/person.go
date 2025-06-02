package service

import (
	"fmt"

	"testovoye/internal/domain"
	"testovoye/internal/repository"
	"testovoye/internal/utils"

	"github.com/google/uuid"
)

type PersonService struct {
	repo          repository.Person
	fetcherExtAPI utils.FetcherAPIByName

	domain string
}

func NewPersonService(repo repository.Person, fetcherExtAPI utils.FetcherAPIByName, domain string) *PersonService {
	return &PersonService{
		repo:          repo,
		fetcherExtAPI: fetcherExtAPI,
		domain:        domain,
	}
}

func (p *PersonService) Add(newPerson NewPersonCreate) (string, error) {
	id := uuid.New()

	age, gender, nationality, err := p.fetcherExtAPI.GetNameMetrics(newPerson.Name)
	if err != nil {
		return "", err
	}

	person := domain.Person{
		ID:          id,
		Name:        newPerson.Name,
		Surname:     newPerson.Surname,
		Patronymic:  newPerson.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	if err := p.repo.Create(person); err != nil {
		return "", fmt.Errorf("failed to create person: %w", err)
	}

	return id.String(), nil
}

func (p *PersonService) GetWithFilter(query domain.PersonQuery) (domain.PaginatedPersons, error) {
	return p.repo.Get(query)
}

func (p *PersonService) DeleteByID(id string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	return p.repo.DeleteByID(id)
}

func (p *PersonService) Edit(editedPerson domain.Person) error {
	if editedPerson.ID.String() == "" {
		return fmt.Errorf("id is required")
	}
	return p.repo.Edit(editedPerson)
}
