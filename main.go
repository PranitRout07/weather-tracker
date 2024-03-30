package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
)

type GetWeatherData struct {
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func query(city string, apiToken string) (GetWeatherData, error) {
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiToken)
	if err != nil {
		return GetWeatherData{}, err
	}
	defer resp.Body.Close()
	var data GetWeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return GetWeatherData{}, err
	}
	data.Main.Kelvin -= 273.15 
	return data, nil
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.TrimPrefix(r.URL.Path, "/weather/")
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
