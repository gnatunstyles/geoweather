package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gnatunstyles/geoweather/cache"
	"github.com/gnatunstyles/geoweather/db"
	"github.com/gnatunstyles/geoweather/internal/config"
	"github.com/gnatunstyles/geoweather/internal/router"
	"github.com/gnatunstyles/geoweather/internal/service"
)

func main() {
	c := &http.Client{}
	myCache := cache.New(10*time.Minute, 10*time.Minute)
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	err = service.InitCities(c, myCache, cfg, database)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(10 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err = service.GetPredictions(c, myCache, cfg, database)
				if err != nil {
					log.Fatal(err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	r := router.New(cfg, database)
	r.Run()

}
