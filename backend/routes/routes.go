package routes

import (
	"api-fio/services/person"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Summary   Получить список Person
// @Tags      person
// @Accept    json
// @Produce   json
// @Param     name        query    string  false  "Имя для фильтра поиска"
// @Param     surname     query    string  false  "Фамилия для фильтра"
// @Param     min_age     query    int     false  "Минимальный возраст"
// @Param     max_age     query    int     false  "Максимальный возраст"
// @Param     gender      query    string  false  "Пол"
// @Param     nationality query    string  false  "Код страны"
// @Param     page        query    int     false  "Номер страницы"
// @Param     page_size   query    int     false  "Размер страницы"
// @Success   200 {object} models.PaginatedPersons
// @Failure   400 {object} models.ErrorResponse
// @Router    /person [get]
func GetPersonHandler(c *gin.Context, logger *zap.Logger, db *gorm.DB) {
	person.GetPerson(c, logger, db)
}

// @Summary   Создать нового Person
// @Tags      person
// @Accept    json
// @Produce   json
// @Param     person  body    models.PersonReqPost    true  "Данные нового Person"
// @Success   200 {object} models.ResponseMessage
// @Failure   400 {object} models.ErrorResponse
// @Failure   500 {object} models.ErrorResponse
// @Router    /person [post]
func AddPersonHandler(c *gin.Context, logger *zap.Logger, db *gorm.DB) {
	person.AddNewPerson(c, logger, db)
}

// @Summary   Удалить Person по ID
// @Tags      person
// @Accept    json
// @Produce   json
// @Param     id   path    int     true  "ID Person"
// @Success   200 {object} models.ResponseMessage
// @Failure   500 {object} models.ErrorResponse
// @Router    /person/{id} [delete]
func DeletePersonHandler(c *gin.Context, logger *zap.Logger, db *gorm.DB) {
	person.DeletePersonByID(c, logger, db)
}

// @Summary   Обновить Person по ID
// @Tags      person
// @Accept    json
// @Produce   json
// @Param     id      path    int            true  "ID Person"
// @Param     person  body    models.PersonReqPut  true  "Новые данные Person"
// @Success   200 {object} models.ResponseMessage
// @Failure   400 {object} models.ErrorResponse
// @Failure   500 {object} models.ErrorResponse
// @Router    /person/{id} [put]
func UpdatePersonHandler(c *gin.Context, logger *zap.Logger, db *gorm.DB) {
	person.UpdatePersonByID(c, logger, db)
}

func InitRoutes(r *gin.Engine, logger *zap.Logger, db *gorm.DB) {
	// Swagger
	r.Static("/docs", "./docs")
	url := ginSwagger.URL("/docs/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Person API
	r.GET("/person", func(c *gin.Context) {
		GetPersonHandler(c, logger, db)
	})
	r.POST("/person", func(c *gin.Context) {
		AddPersonHandler(c, logger, db)
	})
	r.DELETE("/person/:id", func(c *gin.Context) {
		DeletePersonHandler(c, logger, db)
	})
	r.PUT("/person/:id", func(c *gin.Context) {
		UpdatePersonHandler(c, logger, db)
	})
}
