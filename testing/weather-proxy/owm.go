package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OwmCoords struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type OwmMain struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	Humidity float64 `json:"humidity"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}

type OwmWind struct {
	Speed   float32 `json:"speed"`
	Degrees float32 `json:"deg"`
}

type OwmWeather struct {
	Id          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type OwmSys struct {
	Type    int32   `json:"type"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int64   `json:"sunrise"`
	Sunset  int64   `json:"sunset"`
}

type OwmWeatherInfo struct {
	Time       int64        `json:"dt"`
	Coords     OwmCoords    `json:"coord"`
	Weather    []OwmWeather `json:"weather"`
	Base       string       `json:"stations"`
	Main       OwmMain      `json:"main"`
	Visibility float64      `json:"visibility"`
	Wind       OwmWind      `json:"wind"`
	Sys        OwmSys       `json:"sys"`
	Name       string       `json:"name"`
}

func getWeather(zip string, apiKey string) (*OwmWeatherInfo, error) {
	owmUrl := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?zip=%s&APPID=%s", zip, apiKey)
	u, err := url.Parse(owmUrl)
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	jsonBlob, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	w := OwmWeatherInfo{}
	err = json.Unmarshal(jsonBlob, &w)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
