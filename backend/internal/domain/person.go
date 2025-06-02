package domain

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Surname     string    `gorm:"not null"`
	Patronymic  string
	Age         int    `gorm:"not null"`
	Gender      string `gorm:"not null"`
	Nationality string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
