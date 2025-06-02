package repository

import "gorm.io/gorm"

type Repositories struct {
	Person Person
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Person: NewPersonRepository(db),
	}
}
