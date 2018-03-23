package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	fktesting "github.com/fieldkit/cloud/server/api/tool"
	"github.com/fieldkit/cloud/server/backend/ingestion/formatting"
)

type options struct {
	Project         string
	Expedition      string
	DeviceName      string
	WebDevice       bool
	Scheme          string
	Host            string
	RockBlockSerial string
	RockBlock       bool
	CurlOut         bool
	Interval        int
	Username        string
	Password        string
}

func (o *options) getUrl() string {
	return fmt.Sprintf("%s://%s/messages/ingestion", o.Scheme, o.Host)
}

func (o *options) postRawMessage(m *FakeMessage) error {
	if o.CurlOut {
		fmt.Printf("curl -X POST -H 'Content-Type: %s' -d '%s' %s\n", m.ContentType, m.Body, o.getUrl())
	} else {
		body := bytes.NewBufferString(m.Body)
		url := o.getUrl()
		if m.ContentType == formatting.HttpProviderJsonContentType {
			url += "?token=" + "IGNORED"
		}
		_, err := http.Post(url, m.ContentType, body)
		if err != nil {
			log.Fatalf("%s %s", url, err)
		}
	}

	return nil
}

type FakeMessage struct {
	ContentType string
	Body        string
}

func (o *options) createMessage(key string, e *FakeEvent) *FakeMessage {
	if o.RockBlock {
		contentType := "application/x-www-form-urlencoded; charset=UTF-8"
		imei := "300234065114240"
		data := fmt.Sprintf(`%d,BU,LO,3955.00,70.77,%v,%v,933.20,165`, e.Timestamp, e.Coordinates[0], e.Coordinates[1])
		ascii := hex.EncodeToString([]byte(data))

		form := url.Values{}
		form.Set("device_type", "ROCKBLOCK")
		form.Set("serial", o.RockBlockSerial)
		form.Set("momsn", "10")
		form.Set("transmit_time", "17-08-06+12%%3A01%%3A14")
		form.Set("imei", imei)
		form.Set("iridium_cep", "302")
		form.Set("iridium_latitude", "-12.7071")
		form.Set("iridium_longitude", "15.5300")
		form.Set("data", ascii)

		return &FakeMessage{
			ContentType: contentType,
			Body:        form.Encode(),
		}
	}

	if o.WebDevice {
		m := formatting.HttpJsonMessage{
			Location: e.Coordinates,
			Time:     e.Timestamp,
			Device:   key,
			Stream:   "1",
			Values: map[string]interface{}{
				"cpu":      100,
				"temp":     100,
				"humidity": 100,
			},
		}

		bytes, err := json.Marshal(m)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		return &FakeMessage{
			ContentType: formatting.HttpProviderJsonContentType,
			Body:        string(bytes),
		}
	}

	return &FakeMessage{}
}

func main() {
	ctx := context.TODO()

	o := options{}

	flag.StringVar(&o.Project, "project", "www", "project")
	flag.StringVar(&o.Expedition, "expedition", "", "expedition")
	flag.StringVar(&o.DeviceName, "device-name", "test-generator", "device name")
	flag.StringVar(&o.Scheme, "scheme", "http", "scheme to use")
	flag.StringVar(&o.Host, "host", "127.0.0.1:8080", "hostname to use")
	flag.StringVar(&o.Username, "username", "demo-user", "username to use")
	flag.StringVar(&o.Password, "password", "asdfasdfasdf", "password to use")
	flag.BoolVar(&o.WebDevice, "web", false, "fake web device messages")
	flag.StringVar(&o.RockBlockSerial, "rb-serial", "11380", "fake rockblock device serial")
	flag.BoolVar(&o.RockBlock, "rockblock", false, "fake rockblock device messages")
	flag.BoolVar(&o.CurlOut, "curl", false, "write curl commands")
	flag.IntVar(&o.Interval, "interval", 0, "time between messages")

	flag.Parse()

	if !o.WebDevice && !o.RockBlock {
		flag.Usage()
		return
	}

	c, err := fktesting.CreateAndAuthenticate(ctx, o.Host, o.Scheme, o.Username, o.Password)
	if err != nil {
		log.Fatalf("%v", err)
	}

	device, err := fktesting.CreateWebDevice(ctx, c, o.Project, o.DeviceName, "", "1")
	if err != nil {
		log.Fatalf("%v", err)
	}

	path := [][]float64{
		{-114.6846317, 35.4928533},
		{-114.6791397, 35.4967953},
		{-114.6826217, 35.5031313},
		{-114.6789537, 35.5098653},
		{-114.6783727, 35.5155353},
		{-114.6686187, 35.5247973},
	}

	start := time.Now().Add(-24 * time.Hour).Unix()
	end := time.Now().Unix()

	for _, e := range generateFakeEventsAlong(path, 5, start, end) {
		o.postRawMessage(o.createMessage(device.Key, &e))
		time.Sleep(time.Duration(o.Interval) * time.Second)
	}
}

type FakeEvent struct {
	Timestamp   int64
	Coordinates []float64
}

func generateFakeEventsAlong(path [][]float64, numberPerLeg int, startTime, endTime int64) (fakes []FakeEvent) {
	fakes = make([]FakeEvent, 0)

	interval := (endTime - startTime) / int64(len(path)*numberPerLeg)
	now := startTime
	previous := path[0]
	for _, next := range path[1:] {
		for i := 0; i < numberPerLeg; i++ {
			lng := fktesting.MapFloat64(i, 0, numberPerLeg, previous[0], next[0])
			lat := fktesting.MapFloat64(i, 0, numberPerLeg, previous[1], next[1])
			fakes = append(fakes, FakeEvent{
				Timestamp:   now,
				Coordinates: []float64{lng, lat},
			})
			now += interval
		}
		fmt.Printf("%v -> %v\n", previous, next)
		previous = next
	}

	return
}
