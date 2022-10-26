package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sort"

	owm "github.com/briandowns/openweathermap"
	"github.com/gin-gonic/gin"
	"github.com/gnatunstyles/geoweather/internal/config"
	"github.com/gnatunstyles/geoweather/internal/models"
	"github.com/gnatunstyles/geoweather/internal/service"
)

func New(cfg *config.Config, database *sql.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/cities", func(c *gin.Context) {
		sort.Strings(service.CitiesArr)
		c.JSON(http.StatusOK, gin.H{
			"cities": service.CitiesArr,
		})
	})

	r.GET("/shortinfo/:city/:date", func(c *gin.Context) {
		var result []models.Prediction
		city := c.Param("city")
		date := c.Param("date")
		rows, err := database.Query(fmt.Sprintf("SELECT * FROM predictions WHERE city = %s AND date = %s;", city, date))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		defer rows.Close()
		for rows.Next() {
			var pred models.Prediction
			if err := rows.Scan(&pred.City, &pred.Temp, &pred.Date); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}
			result = append(result, pred)
		}
		c.JSON(http.StatusOK, gin.H{
			"predictions": result,
		})
	})

	r.GET("/fullinfo/:city/:date", func(c *gin.Context) {
		city := c.Param("city")
		w, err := owm.NewForecast("5", "F", "FI", cfg.ApiKey) // valid options for first parameter are "5" and "16"
		if err != nil {
			log.Fatalln(err)
		}
		w.DailyByName(
			city,
			5, // five days forecast
		)
		c.JSON(http.StatusOK, gin.H{
			"cities": w,
		})
	})
	return r
}
