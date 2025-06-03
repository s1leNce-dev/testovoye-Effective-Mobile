package v1

type newPersonInput struct {
	// Name имя человека
	// Required: true
	// Min Length: 2
	// Max Length: 100
	Name string `json:"name" binding:"required,min=2,max=100"`
	// Surname фамилия человека
	// Required: true
	// Min Length: 2
	// Max Length: 100
	Surname string `json:"surname" binding:"required,min=2,max=100"`
	// Patronymic отчество (необязательно)
	// Min Length: 2
	// Max Length: 100
	Patronymic string `json:"patronymic,omitempty" binding:"omitempty,min=2,max=100"`
}

type getPersonInput struct {
	// Name фильтрация по имени (partial match)
	// in: query
	// example: Иван
	Name string `form:"name"`
	// Surname фильтрация по фамилии (partial match)
	// in: query
	// example: Петров
	Surname string `form:"surname"`
	// MinAge минимальный возраст (включительно)
	// in: query
	// minimum: 0
	// example: 18
	MinAge int `form:"min_age"`
	// MaxAge максимальный возраст (включительно)
	// in: query
	// minimum: 0
	// example: 65
	MaxAge int `form:"max_age"`
	// Gender фильтрация по полу
	// in: query
	// example: male
	Gender string `form:"gender"`
	// Nationality фильтрация по коду страны (ISO Alpha-2)
	// in: query
	// example: RU
	Nationality string `form:"nationality"`

	// Page номер страницы
	// in: query
	// default: 1
	// minimum: 1
	// example: 1
	Page int `form:"page,default=1"`
	// PageSize количество записей на странице
	// in: query
	// default: 10
	// minimum: 1
	// maximum: 100
	// example: 10
	PageSize int `form:"page_size,default=10"`
}

type updatePersonInput struct {
	// ID идентификатор Person (UUID)
	// Required: true
	// Example: 3fa85f64-5717-4562-b3fc-2c963f66afa6
	ID string `json:"id" binding:"required,uuid4"`
	// Name новое имя (необязательно)
	// Min Length: 2
	// Max Length: 100
	Name string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	// Surname новая фамилия (необязательно)
	// Min Length: 2
	// Max Length: 100
	Surname string `json:"surname,omitempty" binding:"omitempty,min=2,max=100"`
	// Patronymic новое отчество (необязательно)
	// Min Length: 2
	// Max Length: 100
	Patronymic string `json:"patronymic,omitempty" binding:"omitempty,min=2,max=100"`
	// Age новый возраст (необязательно)
	// Minimum: 0
	// Example: 30
	Age int `json:"age,omitempty" binding:"omitempty,min=0"`
	// Gender новый пол (необязательно)
	// Example: female
	Gender string `json:"gender,omitempty"`
	// Nationality новый код страны (необязательно)
	// Example: US
	Nationality string `json:"nationality,omitempty"`
}
