package domain

type PersonQuery struct {
	Name        string `form:"name"`
	Surname     string `form:"surname"`
	Gender      string `form:"gender"`
	Nationality string `form:"nationality"`
	MinAge      int    `form:"min_age"`
	MaxAge      int    `form:"max_age"`

	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

type PaginatedPersons struct {
	Data  []Person
	Total int64

	Page       int
	PageSize   int
	TotalPages int
}
