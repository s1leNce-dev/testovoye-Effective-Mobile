package repository

import (
	"errors"
	"testovoye/internal/domain"

	"gorm.io/gorm"
)

type PersonRepo struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) *PersonRepo {
	return &PersonRepo{db: db}
}

func (r *PersonRepo) Create(newPerson domain.Person) error {
	return r.db.Create(&newPerson).Error
}

func (r *PersonRepo) Get(q domain.PersonQuery) (domain.PaginatedPersons, error) {
	var (
		persons []domain.Person
		total   int64
		result  domain.PaginatedPersons
	)

	tx := r.db.Model(&domain.Person{})

	if q.Name != "" {
		tx = tx.Where("name ILIKE ?", "%"+q.Name+"%")
	}
	if q.Surname != "" {
		tx = tx.Where("surname ILIKE ?", "%"+q.Surname+"%")
	}
	if q.Gender != "" {
		tx = tx.Where("gender = ?", q.Gender)
	}
	if q.Nationality != "" {
		tx = tx.Where("nationality = ?", q.Nationality)
	}
	if q.MinAge > 0 {
		tx = tx.Where("age >= ?", q.MinAge)
	}
	if q.MaxAge > 0 {
		tx = tx.Where("age <= ?", q.MaxAge)
	}

	if err := tx.Count(&total).Error; err != nil {
		return result, err
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 10
	}
	offset := (q.Page - 1) * q.PageSize

	if err := tx.
		Limit(q.PageSize).
		Offset(offset).
		Find(&persons).Error; err != nil {
		return result, err
	}

	totalPages := int((total + int64(q.PageSize) - 1) / int64(q.PageSize))

	result = domain.PaginatedPersons{
		Data:       persons,
		Total:      total,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalPages: totalPages,
	}
	return result, nil
}

func (r *PersonRepo) DeleteByID(id string) error {
	if id == "" {
		return errors.New("id is empty")
	}
	var person domain.Person
	if err := r.db.Where("id = ?", id).First(&person).Error; err != nil {
		return err
	}
	if err := r.db.Delete(person).Error; err != nil {
		return err
	}
	return nil
}

func (r *PersonRepo) Edit(updated domain.Person) error {
	if updated.ID.String() == "" {
		return errors.New("empty id")
	}

	return r.db.Model(&domain.Person{}).
		Where("id = ?", updated.ID).
		Updates(updated).Error
}
