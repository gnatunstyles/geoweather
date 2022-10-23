package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gnatunstyles/geoweather/db"
	"github.com/gnatunstyles/geoweather/internal/config"
	"github.com/gnatunstyles/geoweather/internal/models"
)

const (
	apiGetCity = `https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s`
	apiGetPred = `https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s`
	insertCity = `INSERT INTO cities (city, country, lattitude, longitude) VALUES ($1, $2, $3, $4)`
	insertPred = `INSERT INTO predictions (city, temp, date) VALUES ($1, $2, $3)`
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	citiesArr := []string{"Москва", "Аннино", "Йошкар-Ола", "Тосно", "Краснодар", "Чикаго", "Киев",
		"Берлин", "Вашингтон", "Цюрих", "Париж", "Боровичи", "Барановичи", "Дублин", "Кардифф", "Оттава",
		"Торонто", "Картахена", "Каракас", "Сантьяго"}

	cities := []models.City{}
	c := &http.Client{}
	for _, city := range citiesArr {
		resp, err := c.Get(fmt.Sprintf(apiGetCity, city, cfg.ApiKey))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
		json.Unmarshal(data, &cities)
		fmt.Printf("unpacked in City:\n%#v\n\n", cities[0])
		_, err = database.Exec(insertCity, cities[0].Name, cities[0].Country, cities[0].Lattitude, cities[0].Longitude)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, city := range citiesArr {
		resp, err := c.Get(fmt.Sprintf(apiGetPred, city, cfg.ApiKey))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		p := &models.PredictionParse{}
		// fmt.Println(string(data))
		json.Unmarshal(data, &p)

		for i := range p.List {
			t, err := time.Parse("2006-01-02 15:04:05", p.List[i].Date)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(t)
			if t.Hour() == 12 {
				fmt.Printf("\nTemp:%f\nDate:%s\nCity:%s\n", p.List[i].Main.Temp-273.15, p.List[i].Date, p.City.Name)
				_, err = database.Exec(insertPred, p.City.Name, p.List[i].Main.Temp-273.15, t)
				if err != nil {
					log.Fatal(err)
				}
			}
			// } else {
			// 	fmt.Println(t.Hour())
			// }
		}
	}
	fmt.Println("done")
}
