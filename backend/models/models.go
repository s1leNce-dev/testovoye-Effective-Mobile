package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Surname     string `gorm:"not null" json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `gorm:"not null" json:"age"`
	Gender      string `gorm:"not null" json:"gender"`
	Nationality string `gorm:"not null" json:"nationality"`
}

// ResponseMessage — простой ответ с сообщением
//
// swagger:model ResponseMessage
type ResponseMessage struct {
	// example: success
	Message string `json:"message"`
}

// ErrorResponse — ответ с ошибкой
//
// swagger:model ErrorResponse
type ErrorResponse struct {
	// example: invalid request
	Error string `json:"error"`
}

// PaginatedPersons — страница с данными Person
//
// swagger:model PaginatedPersons
type PaginatedPersons struct {
	// Данные
	Data []Person `json:"data"`
	// Общее число записей
	Total int64 `json:"total"`
	// Текущая страница
	Page int `json:"page"`
	// Размер страницы
	PageSize int `json:"page_size"`
	// Всего страниц
	TotalPages int `json:"total_pages"`
}

// Структура для запроса создания
type PersonReqPost struct {
	Name       string `json:"name"     binding:"required,min=2,max=100"`
	Surname    string `json:"surname"  binding:"required,min=2,max=100"`
	Patronymic string `json:"patronymic,omitempty" binding:"omitempty,min=2,max=100"`
}

type PersonReqPut struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
