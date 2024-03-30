package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/rs/cors"
)

type WeatherAPI struct {
	API string `json:"OpenWeatherAPI"`
}

type WeatherInfo struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

type GetWeatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`

	Weather []WeatherInfo `json:"weather"`
}

func APIConfig(apiToken string) (WeatherAPI, error) {
	var x WeatherAPI
	x.API = apiToken
	fmt.Printf("API token: %+v\n", x.API)
	return x, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!\n"))
}

func query(city string, apiToken string) (GetWeatherData, error) {
	apiConfig, err := APIConfig(apiToken)
	if err != nil {
		return GetWeatherData{}, err
	}
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.API)
	if err != nil {
		return GetWeatherData{}, err
	}
	defer resp.Body.Close()
	var d GetWeatherData
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return GetWeatherData{}, err
	}
	d.Main.Kelvin = d.Main.Kelvin - 273
	d.Main.Kelvin = float64(int(d.Main.Kelvin*100)) / 100
	return d, nil
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			apiToken := os.Getenv("API_TOKEN")
			data, err := query(city, apiToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})


	handler := cors.Default().Handler(mux)

	log.Println("Listening....")
	http.ListenAndServe(":8080", handler)
}
