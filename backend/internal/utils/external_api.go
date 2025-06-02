package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type FetcherAPIByName struct {
	timeout time.Duration
}

func NewFetcherAPIByName(timeout time.Duration) *FetcherAPIByName {
	return &FetcherAPIByName{
		timeout: timeout,
	}
}

func (f *FetcherAPIByName) GetNameMetrics(name string) (age int, gender string, nationality string, err error) {
	client := &http.Client{Timeout: f.timeout}

	ageApi := os.Getenv("API_AGIFY")
	genderApi := os.Getenv("API_GENDERIZE")
	natApi := os.Getenv("API_NATIONALIZE")

	urlsDatas := map[string]interface{}{
		ageApi:    &ageData{},
		genderApi: &genderData{},
		natApi:    &natData{},
	}

	ch := make(chan fetchResult, len(urlsDatas))

	for baseURL, template := range urlsDatas {
		go func(url string, target interface{}) {
			err := getJSON(client, url+name, target)
			if err != nil {
				ch <- fetchResult{URL: url, Err: err}
				return
			}

			ch <- fetchResult{URL: url, Data: target, Err: nil}
		}(baseURL, template)
	}

	var (
		ageVal         float64
		genderVal      string
		nationalityVal string
	)

	for i := 0; i < len(urlsDatas); i++ {
		res := <-ch
		if res.Err != nil {
			return 0, "", "", fmt.Errorf("error fetching %s: %w", res.URL, res.Err)
		}

		log.Println("[DEBUG] res:", res)

		switch res.URL {
		case ageApi:
			ageVal = res.Data.(*ageData).Age
		case genderApi:
			genderVal = res.Data.(*genderData).Gender
		case natApi:
			nat := res.Data.(*natData)
			if len(nat.Country) > 0 {
				nationalityVal = nat.Country[0].CountryID
			}
		}
	}

	if ageVal == 0 || genderVal == "" || nationalityVal == "" {
		return 0, "", "", fmt.Errorf("incomplete data from APIs")
	}

	return int(ageVal), genderVal, nationalityVal, nil
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

type fetchResult struct {
	URL  string
	Data interface{}
	Err  error
}

func getJSON(client *http.Client, url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
