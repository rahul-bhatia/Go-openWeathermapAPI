package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type appiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name       string  `json:"name`
	Visibility float64 `json:"visibility`
	Main       struct {
		Temp float64 `json:"temp`
	} `json:"main"`
}

func loadApiConfig(filename string) (appiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return appiConfigData{}, err
	}
	var c appiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return appiConfigData{}, err
	}
	return c, nil
}
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello world")
}
func query(city string) (weatherData, error) {
	appiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		fmt.Print("err:", err)
		return weatherData{}, err
	}
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + appiConfig.OpenWeatherMapApiKey)
	fmt.Println("Response ->", resp.Body)
	if err != nil {
		fmt.Print("err:", err)
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		fmt.Print("err:", err)
		return weatherData{}, err
	}
	return d, nil
}
func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})
	http.ListenAndServe(":8080", nil)
}
