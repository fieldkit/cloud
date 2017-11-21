package testing

import (
	"context"
	"fmt"
	fk "github.com/fieldkit/cloud/server/api/client"
	goaclient "github.com/goadesign/goa/client"
	"log"
	"net/http"
)

func CreateWebDevice(ctx context.Context, c *fk.Client, projectSlug, deviceName string) (d *fk.DeviceInput, err error) {
	res, err := c.ListProject(ctx, fk.ListProjectPath())
	if err != nil {
		log.Fatalf("%v", err)
	}
	projects, err := c.DecodeProjects(res)
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, project := range projects.Projects {
		log.Printf("Project: %+v\n", *project)

		if project.Slug == projectSlug {
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

					if device.Name == deviceName {
						theDevice = device
					}
				}

				if theDevice == nil {
					log.Printf("Creating new Device %v", deviceName)

					addPayload := fk.AddDeviceInputPayload{
						Name: deviceName,
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

				return theDevice, nil
			}
		}
	}

	return nil, fmt.Errorf("Unable to create Web device")
}

func MapInt64(v, minV, maxV int, minRange, maxRange int64) int64 {
	return int64(MapFloat64(v, minV, maxV, float64(minRange), float64(maxRange)))
}

func MapFloat64(v, minV, maxV int, minRange, maxRange float64) float64 {
	if minV == maxV {
		return minRange
	}
	return float64(v-minV)/float64(maxV-minV)*(maxRange-minRange) + minRange
}

func CreateAndAuthenticate(ctx context.Context, host, scheme, username, password string) (*fk.Client, error) {
	httpClient := newHTTPClient()
	c := fk.New(goaclient.HTTPClientDoer(httpClient))
	c.Host = host
	c.Scheme = scheme

	loginPayload := fk.LoginPayload{}
	loginPayload.Username = username
	loginPayload.Password = password
	res, err := c.LoginUser(ctx, fk.LoginUserPath(), &loginPayload)
	if err != nil {
		return nil, err
	}

	key := res.Header.Get("Authorization")
	jwtSigner := newJWTSigner(key, "%s")
	c.SetJWTSigner(jwtSigner)

	return c, nil
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
