package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fieldkit/cloud/server/backend/ingestion/formatting"
	fktesting "github.com/fieldkit/cloud/server/tools"
)

type options struct {
	Project    string
	Expedition string
	DeviceName string
	ApiKey     string
	Zip        string
	Host       string
	Scheme     string
	Username   string
	Password   string
}

func main() {
	ctx := context.TODO()

	o := options{}

	flag.StringVar(&o.Project, "project", "www", "project")
	flag.StringVar(&o.DeviceName, "device-name", "weather-proxy", "device name")
	flag.StringVar(&o.ApiKey, "api-key", "", "owm api key")
	flag.StringVar(&o.Zip, "zip", "", "zip")
	flag.StringVar(&o.Scheme, "scheme", "http", "scheme to use")
	flag.StringVar(&o.Host, "host", "127.0.0.1:8080", "hostname to use")
	flag.StringVar(&o.Username, "username", "demo-user", "username to use")
	flag.StringVar(&o.Password, "password", "asdfasdfasdf", "password to use")

	flag.Parse()

	if o.ApiKey == "" || o.Zip == "" {
		flag.Usage()
		os.Exit(2)
	}

	w, err := getWeather(o.Zip, o.ApiKey)
	if err != nil {
		log.Fatalf("Unable to get Weather: %v", err)
	}
	fmt.Printf("%+v", w)

	c, err := fktesting.CreateAndAuthenticate(ctx, o.Host, o.Scheme, o.Username, o.Password)
	if err != nil {
		log.Fatalf("%v", err)
	}

	device, err := fktesting.CreateWebDevice(ctx, c, o.Project, o.DeviceName, "", "1")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%v\n", device)

	m := createFieldKitMessageFromWeather(w, device.Key)
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	body := bytes.NewBufferString(string(b))
	url := fmt.Sprintf("%s://%s/messages/ingestion", o.Scheme, o.Host)
	url += "?token=" + "IGNORED"
	_, err = http.Post(url, formatting.HttpProviderJsonContentType, body)
	if err != nil {
		log.Fatalf("%s %s", url, err)
	}
}

func kToF(k float64) float64 {
	return k*9/5.0 - 459.67
}

func createFieldKitMessageFromWeather(w *OwmWeatherInfo, device string) *formatting.HttpJsonMessage {
	return &formatting.HttpJsonMessage{
		Location: []float64{w.Coords.Longitude, w.Coords.Latitude},
		Time:     w.Time,
		Device:   device,
		Stream:   "1",
		Values: map[string]interface{}{
			"temp":       kToF(w.Main.Temp),
			"pressure":   w.Main.Pressure, // hPa
			"humidity":   w.Main.Humidity, // %
			"visibility": w.Visibility,    // meters
		},
	}
}
