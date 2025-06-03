package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

type FetcherAPIByName struct {
	timeout time.Duration
}

func NewFetcherAPIByName(timeout time.Duration) *FetcherAPIByName {
	return &FetcherAPIByName{timeout: timeout}
}

func (f *FetcherAPIByName) GetNameMetrics(name string) (age int, gender string, nationality string, err error) {
	client := &http.Client{Timeout: f.timeout}

	ageURL := os.Getenv("API_AGIFY")
	genderURL := os.Getenv("API_GENDERIZE")
	natURL := os.Getenv("API_NATIONALIZE")

	var (
		ageVal         float64
		genderVal      string
		nationalityVal string
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		var data ageData
		if err := getJSON(client, ageURL+name, &data); err != nil {
			return fmt.Errorf("error fetching age from %s: %w", ageURL, err)
		}
		ageVal = data.Age
		log.Printf("[DEBUG] fetched age: %v\n", ageVal)
		return nil
	})

	eg.Go(func() error {
		var data genderData
		if err := getJSON(client, genderURL+name, &data); err != nil {
			return fmt.Errorf("error fetching gender from %s: %w", genderURL, err)
		}
		genderVal = data.Gender
		log.Printf("[DEBUG] fetched gender: %s\n", genderVal)
		return nil
	})

	eg.Go(func() error {
		var data natData
		if err := getJSON(client, natURL+name, &data); err != nil {
			return fmt.Errorf("error fetching nationality from %s: %w", natURL, err)
		}
		if len(data.Country) > 0 {
			nationalityVal = data.Country[0].CountryID
		}
		log.Printf("[DEBUG] fetched nationality: %s\n", nationalityVal)
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 0, "", "", err
	}

	if ageVal == 0 || genderVal == "" || nationalityVal == "" {
		return 0, "", "", fmt.Errorf("incomplete data from APIs: age=%v, gender=%q, nationality=%q", ageVal, genderVal, nationalityVal)
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

func getJSON(client *http.Client, url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
