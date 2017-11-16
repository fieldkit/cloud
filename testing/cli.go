package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	fk "github.com/fieldkit/cloud/server/api/client"
	"github.com/fieldkit/cloud/server/backend/ingestion"
	goaclient "github.com/goadesign/goa/client"
	"log"
	"net/http"
	"time"
)

type options struct {
	Project    string
	Expedition string
	Name       string
	Scheme     string
	Host       string
	CurlOut    bool
}

func (o *options) getUrl() string {
	return fmt.Sprintf("%s://%s/messages/ingestion", o.Scheme, o.Host)
}

func (o *options) postRawMessage(contentType string, raw string) error {
	body := bytes.NewBufferString(raw)
	url := o.getUrl()
	if contentType == ingestion.HttpProviderJsonContentType {
		url += "?token=" + "IGNORED"
	}
	_, err := http.Post(url, contentType, body)
	if err != nil {
		log.Fatalf("%s %s", url, err)
	}
	return nil
}

func generateRockBlock() (raw []byte, err error) {
	return nil, nil
}

func main() {
	ctx := context.TODO()

	o := options{}

	flag.StringVar(&o.Project, "project", "www", "project")
	flag.StringVar(&o.Expedition, "expedition", "", "expedition")
	flag.StringVar(&o.Name, "name", "test-generator", "device name")
	flag.StringVar(&o.Scheme, "scheme", "http", "scheme to use")
	flag.StringVar(&o.Host, "host", "127.0.0.1:8080", "hostname to use")
	flag.BoolVar(&o.CurlOut, "curl", false, "write curl commands")

	flag.Parse()

	if true {
		httpClient := newHTTPClient()
		c := fk.New(goaclient.HTTPClientDoer(httpClient))
		c.Scheme = "http"
		c.Host = "127.0.0.1:8080"

		loginPayload := fk.LoginPayload{}
		loginPayload.Username = "demo-user"
		loginPayload.Password = "asdfasdfasdf"
		res, err := c.LoginUser(ctx, fk.LoginUserPath(), &loginPayload)
		if err != nil {
			log.Fatalf("%v", err)
		}

		key := res.Header.Get("Authorization")
		jwtSigner := newJWTSigner(key, "%s")
		c.SetJWTSigner(jwtSigner)

		res, err = c.ListProject(ctx, fk.ListProjectPath())
		if err != nil {
			log.Fatalf("%v", err)
		}
		projects, err := c.DecodeProjects(res)
		if err != nil {
			log.Fatalf("%v", err)
		}

		for _, project := range projects.Projects {
			log.Printf("Project: %+v\n", *project)

			if project.Slug == o.Project {
				res, err = c.ListExpedition(ctx, fk.ListExpeditionPath(project.Slug))
				if err != nil {
					log.Fatalf("%v", err)
				}
				expeditions, err := c.DecodeExpeditions(res)
				if err != nil {
					log.Fatalf("%v", err)
				}

				for _, exp := range expeditions.Expeditions {
					log.Printf("Expedition: %+v\n", *exp)

					res, err = c.ListDevice(ctx, fk.ListDevicePath(project.Slug, exp.Slug))
					if err != nil {
						log.Fatalf("%v", err)
					}
					devices, err := c.DecodeDeviceInputs(res)
					if err != nil {
						log.Fatalf("%v", err)
					}

					var theDevice *fk.DeviceInput
					for _, device := range devices.DeviceInputs {
						log.Printf("Device: %+v\n", *device)

						if device.Name == o.Name {
							theDevice = device
						}
					}

					if theDevice == nil {
						log.Printf("Creating new Device %v", o.Name)

						addPayload := fk.AddDeviceInputPayload{
							Key:  "HTTP-" + o.Name,
							Name: o.Name,
						}
						res, err = c.AddDevice(ctx, fk.AddDevicePath(exp.ID), &addPayload)
						if err != nil {
							log.Fatalf("%v", err)
						}

						added, err := c.DecodeDeviceInput(res)
						if err != nil {
							log.Fatalf("%v", err)
						}

						theDevice = added
					}

					schema := fk.UpdateDeviceInputSchemaPayload{
						Active:     true,
						Key:        "1",
						JSONSchema: `{ "UseProviderTime": true, "UseProviderLocation": true }`,
					}
					_, err = c.UpdateSchemaDevice(ctx, fk.UpdateSchemaDevicePath(theDevice.ID), &schema)
					if err != nil {
						log.Fatalf("%v", err)
					}

					path := [][]float64{
						{35.4928533, -114.6846317},
						{35.4967953, -114.6791397},
						{35.5031313, -114.6826217},
						{35.5098653, -114.6789537},
						{35.5155353, -114.6783727},
						{35.5247973, -114.6686187},
					}

					start := time.Now().Add(-24 * time.Hour).Unix()
					end := time.Now().Unix()

					for _, e := range generateFakeEventsAlong(path, 5, start, end) {
						m := ingestion.HttpJsonMessage{
							Location: e.Coordinates,
							Time:     e.Timestamp,
							Device:   o.Name,
							Stream:   "1",
							Values: map[string]string{
								"cpu":      "100",
								"temp":     "100",
								"humidity": "100",
							},
						}

						bytes, err := json.Marshal(m)
						if err != nil {
							log.Fatalf("Error %v", err)
						}

						log.Printf("%s", string(bytes))

						o.postRawMessage(ingestion.HttpProviderJsonContentType, string(bytes))
					}
				}
			}
		}

		return
	}

	if false {
		path := [][]float64{
			{35.4928533, -114.6846317},
			{35.4967953, -114.6791397},
			{35.5031313, -114.6826217},
			{35.5098653, -114.6789537},
			{35.5155353, -114.6783727},
			{35.5247973, -114.6686187},
		}

		start := time.Now().Add(-24 * time.Hour).Unix()
		end := time.Now().Unix()

		fakeEvents := generateFakeEventsAlong(path, 5, start, end)

		contentType := "application/x-www-form-urlencoded; charset=UTF-8"
		serial := "11380"
		imei := "300234065114240"

		for _, e := range fakeEvents {
			data := fmt.Sprintf(`%d,BU,LO,3955.00,70.77,%v,%v,933.20,165`, e.Timestamp, e.Coordinates[0], e.Coordinates[1])
			ascii := hex.EncodeToString([]byte(data))
			raw := fmt.Sprintf("device_type=ROCKBLOCK&serial=%s&momsn=1086&transmit_time=17-08-06+12%%3A01%%3A14&imei=%s&iridium_latitude=-12.7071&iridium_longitude=15.5300&iridium_cep=302&data=%s", serial, imei, ascii)

			if o.CurlOut {
				fmt.Printf("curl -X POST -H 'Content-Type: %s' -d '%s' %s\n", contentType, raw, o.getUrl())
			} else {
				fmt.Printf("%s\n", data)
				o.postRawMessage(contentType, raw)
			}
		}
	}
}

type FakeEvent struct {
	Timestamp   int64
	Coordinates []float64
}

func mapInt64(v, minV, maxV int, minRange, maxRange int64) int64 {
	return int64(mapFloat64(v, minV, maxV, float64(minRange), float64(maxRange)))
}

func mapFloat64(v, minV, maxV int, minRange, maxRange float64) float64 {
	if minV == maxV {
		return minRange
	}
	return float64(v-minV)/float64(maxV-minV)*(maxRange-minRange) + minRange
}

func generateFakeEventsAlong(path [][]float64, numberPerLeg int, startTime, endTime int64) (fakes []FakeEvent) {
	fakes = make([]FakeEvent, 0)

	interval := (endTime - startTime) / int64(len(path)*numberPerLeg)
	now := startTime
	previous := path[0]
	for _, next := range path[1:] {
		for i := 0; i < numberPerLeg; i++ {
			lat := mapFloat64(i, 0, numberPerLeg, previous[0], next[0])
			lng := mapFloat64(i, 0, numberPerLeg, previous[1], next[1])
			fakes = append(fakes, FakeEvent{
				Timestamp:   now,
				Coordinates: []float64{lat, lng},
			})
			now += interval
		}
		fmt.Printf("%v -> %v\n", previous, next)
		previous = next
	}

	return
}

func newHTTPClient() *http.Client {
	return http.DefaultClient
}

func newJWTSigner(key, format string) goaclient.Signer {
	return &goaclient.APIKeySigner{
		SignQuery: false,
		KeyName:   "Authorization",
		KeyValue:  key,
		Format:    format,
	}

}
