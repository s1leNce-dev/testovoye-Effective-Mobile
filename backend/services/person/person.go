package person

import (
	"api-fio/models"
	"api-fio/utils/snippets"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func getJSON(client *http.Client, url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

type ageData struct {
	Age float64 `json:"age"`
}
type genderData struct {
	Gender string `json:"gender"`
}
type natData struct {
	Country []struct {
		CountryID string `json:"country_id"`
	} `json:"country"`
}

var urlsDatas = map[string]interface{}{
	"https://api.agify.io/?name=":       &ageData{},
	"https://api.genderize.io/?name=":   &genderData{},
	"https://api.nationalize.io/?name=": &natData{},
}

type fetchResult struct {
	URL  string
	Data interface{}
	Err  error
}

func fetchAll(name string, logger *zap.Logger) ([]fetchResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	ch := make(chan fetchResult, len(urlsDatas))

	for baseURL, template := range urlsDatas {
		go func(url string, target interface{}) {
			err := getJSON(client, url+name, target)
			snippets.HandleDebugLogs(logger, fmt.Sprint("target:", target))
			ch <- fetchResult{URL: url, Data: target, Err: err}
		}(baseURL, template)
	}

	var results []fetchResult
	for i := 0; i < len(urlsDatas); i++ {
		res := <-ch
		if res.Err != nil {
			return nil, res.Err
		}
		results = append(results, res)
	}

	snippets.HandleDebugLogs(logger, "fetched")

	return results, nil
}

func AddNewPerson(c *gin.Context, logger *zap.Logger, db *gorm.DB) {
	var req struct {
		Name       string `json:"name" binding:"required,min=2,max=100"`
		Surname    string `json:"surname" binding:"required,min=2,max=100"`
		Patronymic string `json:"patronymic" binding:"omitempty,min=2,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		snippets.HandleErrorJSONAnswer(c, logger, 400, err.Error(), err.Error(), "[JSON]")
		return
	}

	results, err := fetchAll(req.Name, logger)
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, logger, 500, err.Error(), err.Error(), "[API]")
		return
	}

	var age ageData
	var gender genderData
	var nat natData

	for _, r := range results {
		switch r.URL {
		case "https://api.agify.io/?name=":
			age = *r.Data.(*ageData)
		case "https://api.genderize.io/?name=":
			gender = *r.Data.(*genderData)
		case "https://api.nationalize.io/?name=":
			nat = *r.Data.(*natData)
		}
	}

	if age.Age == 0 || gender.Gender == "" || nat.Country == nil {
		snippets.HandleDebugLogs(logger, fmt.Sprint("After failed fetch: ", age, gender, nat))
		snippets.HandleErrorJSONAnswer(c, logger, 500, "failed to fetch", "failed to fetch", "[API]")
		return
	}

	newPerson := models.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         int(age.Age),
		Gender:      gender.Gender,
		Nationality: nat.Country[0].CountryID,
	}

	if err := db.Create(&newPerson).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, logger, 500, "failed to create a person", err.Error(), "[PGSQL]")
		return
	}

	snippets.HandleInfoLogs(logger, "New person was added", "[PERSON] [POST]")

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// get people
type personQuery struct {
	Name        string `form:"name"`
	Surname     string `form:"surname"`
	MinAge      int    `form:"min_age"`
	MaxAge      int    `form:"max_age"`
	Gender      string `form:"gender"`
	Nationality string `form:"nationality"`

	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

func GetPerson(c *gin.Context, logger *zap.Logger, db *gorm.DB) {
	var q personQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		snippets.HandleErrorJSONAnswer(c, logger, 400, err.Error(), err.Error(), "[JSON]")
		return
	}

	tx := db.Model(&models.Person{})

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

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, logger, 500, err.Error(), err.Error(), "[PGSQL]")
		return
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 10
	}
	offset := (q.Page - 1) * q.PageSize

	var persons []models.Person
	if err := tx.
		Limit(q.PageSize).
		Offset(offset).
		Find(&persons).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, logger, 500, err.Error(), err.Error(), "[PGSQL]")
		return
	}

	totalPages := int((total + int64(q.PageSize) - 1) / int64(q.PageSize))
	c.JSON(http.StatusOK, models.PaginatedPersons{
		Data:       persons,
		Total:      total,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalPages: totalPages,
	})
}

// delete
func DeletePersonByID(c *gin.Context, logger *zap.Logger, db *gorm.DB) {
	id := c.Param("id")

	if err := db.Delete(&models.Person{}, id).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, logger, 500, "Failed to delete a person", err.Error(), "[PGSQL]")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// update
func UpdatePersonByID(c *gin.Context, logger *zap.Logger, db *gorm.DB) {
	id := c.Param("id")

	var updatedPerson models.PersonReqPut
	if err := c.ShouldBindJSON(&updatedPerson); err != nil {
		snippets.HandleErrorJSONAnswer(c, logger, 400, err.Error(), err.Error(), "[JSON]")
		return
	}

	if err := db.Model(&models.Person{}).Where("id = ?", id).Updates(updatedPerson).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, logger, 500, "failed to delete a person", err.Error(), "[PGSQL]")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
