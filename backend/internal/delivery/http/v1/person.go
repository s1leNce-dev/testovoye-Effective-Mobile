package v1

import (
	"net/http"
	"testovoye/internal/domain"
	"testovoye/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *HandlerV1) initPersonRoutes(api *gin.RouterGroup) {
	person := api.Group("/person")
	{
		person.GET("/", h.getPersonData)
		person.POST("/", h.createNewPerson)
		person.PUT("/", h.updatePerson)
		person.DELETE("/", h.deletePerson)
	}
}

// @Summary Создать нового Person
// @Tags Person
// @Accept json
// @Produce json
// @Param input body newPersonInput true "Данные нового Person"
// @Success 200 {object} createPersonResponse
// @Failure 400 {object} ErrorResponse "Неверный формат входных данных"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /person [post]
func (h *HandlerV1) createNewPerson(c *gin.Context) {
	var input newPersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body", err.Error())
		return
	}

	id, err := h.services.Person.Add(service.NewPersonCreate{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: input.Patronymic,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to create", err.Error())
		return
	}

	c.JSON(http.StatusOK, createPersonResponse{ID: id})
}

// @Summary Получить список Person с фильтрацией и пагинацией
// @Tags Person
// @Accept json
// @Produce json
// @Param name query string false "Фильтрация по имени (partial match)"
// @Param surname query string false "Фильтрация по фамилии (partial match)"
// @Param min_age query int false "Минимальный возраст" minimum(0)
// @Param max_age query int false "Максимальный возраст" minimum(0)
// @Param gender query string false "Пол (male/female)"
// @Param nationality query string false "Код страны (ISO Alpha-2)"
// @Param page query int false "Страница" default(1) minimum(1)
// @Param page_size query int false "Размер страницы" default(10) minimum(1) maximum(100)
// @Success 200 {object} domain.PaginatedPersons
// @Failure 400 {object} ErrorResponse "Неверные query-параметры"
// @Failure 500 {object} ErrorResponse "Ошибка при выборке из БД"
// @Router /person [get]
func (h *HandlerV1) getPersonData(c *gin.Context) {
	var queryInput getPersonInput
	if err := c.ShouldBindQuery(&queryInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid query parameters", err.Error())
		return
	}

	svcQuery := domain.PersonQuery{
		Name:        queryInput.Name,
		Surname:     queryInput.Surname,
		MinAge:      queryInput.MinAge,
		MaxAge:      queryInput.MaxAge,
		Gender:      queryInput.Gender,
		Nationality: queryInput.Nationality,
		Page:        queryInput.Page,
		PageSize:    queryInput.PageSize,
	}

	result, err := h.services.Person.GetWithFilter(svcQuery)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to fetch persons", err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Обновить существующего Person
// @Tags Person
// @Accept json
// @Produce json
// @Param input body updatePersonInput true "Данные для обновления Person"
// @Success 200 {object} simpleMessageResponse
// @Failure 400 {object} ErrorResponse "Неправильный формат входных данных или неверный UUID"
// @Failure 500 {object} ErrorResponse "Ошибка при обновлении в БД"
// @Router /person [put]
func (h *HandlerV1) updatePerson(c *gin.Context) {
	var input updatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body", err.Error())
		return
	}

	id, err := uuid.Parse(input.ID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid type of id", err.Error())
		return
	}

	svcUpdate := domain.Person{
		ID:          id,
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  input.Patronymic,
		Age:         input.Age,
		Gender:      input.Gender,
		Nationality: input.Nationality,
	}

	if err := h.services.Person.Edit(svcUpdate); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to update person", err.Error())
		return
	}

	c.JSON(http.StatusOK, simpleMessageResponse{Message: "success"})
}

// @Summary Удалить Person по ID
// @Tags Person
// @Accept json
// @Produce json
// @Param id query string true "UUID идентификатор Person"
// @Success 200 {object} simpleMessageResponse
// @Failure 400 {object} ErrorResponse "ID не указан или неверный формат"
// @Failure 500 {object} ErrorResponse "Ошибка при удалении из БД"
// @Router /person [delete]
func (h *HandlerV1) deletePerson(c *gin.Context) {
	var input deletePersonInput
	if err := c.ShouldBindQuery(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "id is required", err.Error())
		return
	}

	if err := h.services.Person.DeleteByID(input.ID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to delete person", err.Error())
		return
	}

	c.JSON(http.StatusOK, simpleMessageResponse{Message: "success"})
}
