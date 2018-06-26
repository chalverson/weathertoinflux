package main

import (
	owm "github.com/briandowns/openweathermap"
	"log"
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"github.com/spf13/viper"
	"fmt"
	"strconv"
	"os/user"
	"path/filepath"
)

func main() {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	viper.SetConfigName("weathertoinflux")
	viper.AddConfigPath(filepath.Join(homeDir, ".weathertoinflux"))
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     viper.GetString("db.address"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
	})

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: viper.GetString("db.name"),
	})

	if err != nil {
		log.Fatal(err)
	}

	cities := make([]map[string]string, 0)
	var m map[string]string
	citiesI := viper.Get("openweather.cities")
	citiesSES := citiesI.([]interface{})
	for _, cty := range citiesSES {
		citiesMap := cty.(map[interface{}]interface{})
		m = make(map[string]string)
		for k, v := range citiesMap {
			m[k.(string)] = v.(string)
		}
		cities = append(cities, m)
	}

	for _, city := range cities {
		w, err := owm.NewCurrent("F", "en", viper.GetString("openweather.apikey"))
		if err != nil {
			log.Fatal(err)
		}

		if city["type"] == "id" {
			cityId, err := strconv.Atoi(city["id"])
			if err != nil {
				log.Printf("Could not parse city id int: %v\n", city["id"])
				continue
			}
			w.CurrentByID(cityId)
		}

		if city["type"] == "latlon" {
			lat, err := strconv.ParseFloat(city["lat"], 64)
			if err != nil {
				log.Printf("Could not parse lat float: %v\n", city["lat"])
				continue
			}
			lon, err := strconv.ParseFloat(city["lon"], 64)
			if err != nil {
				log.Printf("Could not parse lon float: %v\n", city["lon"])
				continue
			}

			fmt.Printf("Lat: %v Lon: %v\n", lat, lon)
			err = w.CurrentByCoordinates(&owm.Coordinates{
				Latitude:  lat,
				Longitude: lon,
			})
			if err != nil {
				log.Printf("Error in ByCoords: %v", err)
			}
		}

		if city["type"] == "zip" {
			zip, err := strconv.Atoi(city["zip"])
			if err != nil {
				log.Printf("Could not parse zip int: %v\n", city["zip"])
				continue
			}
			w.CurrentByZip(zip, city["countryCode"])
		}

		tagsTemperature := map[string]string{"deviceName": city["label"], "unit": "F"}
		fieldsTemperature := map[string]interface{}{
			"value": w.Main.Temp,
		}
		pt, err := client.NewPoint("temperature", tagsTemperature, fieldsTemperature, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		tagsHumidity := map[string]string{"deviceName": city["label"], "unit": "%"}
		fieldsHumidity := map[string]interface{}{
			"value": float64(w.Main.Humidity),
		}
		ptHumidity, err := client.NewPoint("humidity", tagsHumidity, fieldsHumidity, time.Now())
		bp.AddPoint(ptHumidity)
	}

	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}
