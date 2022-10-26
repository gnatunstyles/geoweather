package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gnatunstyles/geoweather/cache"
	"github.com/gnatunstyles/geoweather/internal/config"
	"github.com/gnatunstyles/geoweather/internal/models"
)

const (
	apiGetCity = `https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s`
	apiGetPred = `https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s`
	insertCity = `INSERT INTO cities (city, country, lattitude, longitude) VALUES ($1, $2, $3, $4)`
	insertPred = `INSERT INTO predictions (city, temp, date) VALUES ($1, $2, $3);`
)

var (
	CitiesArr = []string{"Moscow", "Annino", "Yoshkar-Ola", "Tosno", "Krasnodar", "Chicago", "Kyiv",
		"Berlin", "Washington", "Zurich", "Paris", "Borovichi", "Baranavichy", "Dublin", "Cardiff", "Ottawa",
		"Toronto", "Cartagena", "Caracas", "Irkutsk"}
	savedCities = []string{}
)

func InitCities(c *http.Client, cache *cache.Cache, cfg *config.Config, db *sql.DB) error {
	cities := []models.City{}
	for _, city := range CitiesArr {
		resp, err := c.Get(fmt.Sprintf(apiGetCity, city, cfg.ApiKey))
		if err != nil {
			return fmt.Errorf("error during getting the response:%e", err)
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error during reading the response:%e", err)
		}
		err = json.Unmarshal(data, &cities)
		if err != nil {
			return fmt.Errorf("error during unmarshaling:%e", err)
		}
		_, err = db.Exec(insertCity, cities[0].Name, cities[0].Country, cities[0].Lattitude, cities[0].Longitude)
		if err != nil {
			return fmt.Errorf("error during saving city to db:%e", err)
		}
		savedCities = append(savedCities, cities[0].Name)
		cache.Set(cities[0].Name, cities[0], 10*time.Minute)
		log.Printf("%s city info was saved to db and cached successfully.\n", cities[0].Name)
	}
	return nil
}

func GetPredictions(c *http.Client, cache *cache.Cache, cfg *config.Config, db *sql.DB) error {
	fmt.Println(cache)
	for _, city := range savedCities {
		cityInfo, _ := cache.Get(city)
		fmt.Println(cityInfo)
		resp, err := c.Get(fmt.Sprintf(apiGetPred, cityInfo.Lattitude, cityInfo.Longitude, cfg.ApiKey))
		if err != nil {
			return fmt.Errorf("error during getting the response:%e", err)
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error during getting data from the response:%e", err)
		}
		p := &models.PredictionParse{}
		err = json.Unmarshal(data, &p)
		if err != nil {
			return fmt.Errorf("error during unmarshaling:%e", err)
		}
		for i := range p.List {
			t, err := time.Parse("2006-01-02 15:04:05", p.List[i].Date)
			if err != nil {
				return fmt.Errorf("error during parsing the time:%e", err)
			}
			if t.Hour() == 12 {
				_, err = db.Exec(insertPred, city, p.List[i].Main.Temp, t)
				if err != nil {
					return fmt.Errorf("error during saving prediction to db:%e", err)
				}
			}
		}
		log.Printf("predictions for the next 5 days for the %s city were saved to db successfully.\n", city)
	}
	return nil
}
