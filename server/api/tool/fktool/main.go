package main

import (
	"context"
	"flag"
	fktesting "github.com/fieldkit/cloud/server/api/tool"
	"log"
)

type options struct {
	Scheme   string
	Host     string
	Username string
	Password string

	Project    string
	Expedition string

	DeviceName string
	DeviceId   string
	StreamName string
}

func main() {
	ctx := context.TODO()

	o := options{}

	flag.StringVar(&o.Scheme, "scheme", "http", "fk instance scheme")
	flag.StringVar(&o.Host, "host", "127.0.0.1:8080", "fk instance hostname")
	flag.StringVar(&o.Username, "username", "demo-user", "username")
	flag.StringVar(&o.Password, "password", "asdfasdfasdf", "password")

	flag.StringVar(&o.Project, "project", "www", "project")
	flag.StringVar(&o.Expedition, "expedition", "", "expedition")
	flag.StringVar(&o.DeviceName, "device-name", "", "device name")
	flag.StringVar(&o.DeviceId, "device-id", "", "device id")
	flag.StringVar(&o.StreamName, "stream-name", "", "stream name")

	flag.Parse()

	c, err := fktesting.CreateAndAuthenticate(ctx, o.Host, o.Scheme, o.Username, o.Password)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("Authenticated as %s (%s)", o.Username, o.Host)

	if o.DeviceName != "" {
		device, err := fktesting.CreateWebDevice(ctx, c, o.Project, o.DeviceName, o.DeviceId, "")
		if err != nil {
			log.Fatalf("Error creating web device: %v", err)
		}

		log.Printf("Associated %v", device)
	}
}
