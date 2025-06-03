package domain

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"type:varchar(100);not null;index:idx_person_name"`
	Surname     string    `gorm:"type:varchar(100);not null;index:idx_person_surname"`
	Patronymic  string    `gorm:"type:varchar(100)"`
	Age         int       `gorm:"not null;index:idx_person_age"`
	Gender      string    `gorm:"type:varchar(10);not null;index:idx_person_gender"`
	Nationality string    `gorm:"type:char(2);not null;index:idx_person_nationality"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
