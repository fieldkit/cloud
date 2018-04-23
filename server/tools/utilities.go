package utilities

import (
	"context"
	"fmt"
	fk "github.com/fieldkit/cloud/server/api/client"
	goaclient "github.com/goadesign/goa/client"
	"log"
	"net/http"
)

func UpdateLocation(ctx context.Context, c *fk.Client, device *fk.DeviceSource, longitude, latitude float64) error {
	updateLocationPayload := fk.UpdateDeviceSourceLocationPayload{
		Longitude: longitude,
		Latitude:  latitude,
	}
	res, err := c.UpdateLocationDevice(ctx, fk.UpdateLocationDevicePath(device.ID), &updateLocationPayload)
	if err != nil {
		log.Fatalf("Error updating device location: %v", err)
	}

	updated, err := c.DecodeDeviceSource(res)
	if err != nil {
		log.Fatalf("Error adding device location: %v", err)
	}

	log.Printf("Updated: %+v", updated)

	return nil
}

func CreateWebDevice(ctx context.Context, c *fk.Client, projectSlug, deviceName, deviceId, streamName string) (d *fk.DeviceSource, err error) {
	res, err := c.ListProject(ctx, fk.ListProjectPath())
	if err != nil {
		log.Fatalf("Error listing projects: %v", err)
	}
	projects, err := c.DecodeProjects(res)
	if err != nil {
		log.Fatalf("Error listing projects: %v", err)
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
				devices, err := c.DecodeDeviceSources(res)
				if err != nil {
					log.Fatalf("%v", err)
				}

				var theDevice *fk.DeviceSource
				for _, device := range devices.DeviceSources {
					log.Printf("Device: %+v\n", *device)

					if device.Name == deviceName {
						theDevice = device
					}
				}

				if theDevice == nil {
					log.Printf("Creating new Device %v", deviceName)

					addPayload := fk.AddDeviceSourcePayload{
						Name: deviceName,
						Key:  deviceId,
					}
					res, err = c.AddDevice(ctx, fk.AddDevicePath(exp.ID), &addPayload)
					if err != nil {
						log.Fatalf("Error adding device: %v", err)
					}

					added, err := c.DecodeDeviceSource(res)
					if err != nil {
						log.Fatalf("Error adding device: %v", err)
					}

					theDevice = added
				} else {
					if deviceId != "" && theDevice.Key != deviceId {
						log.Printf("Device id/key is different (%s != %s) fixing...", deviceId, theDevice.Key)

						updatePayload := fk.UpdateDeviceSourcePayload{
							Name: deviceName,
							Key:  deviceId,
						}
						res, err = c.UpdateDevice(ctx, fk.UpdateDevicePath(theDevice.ID), &updatePayload)
						if err != nil {
							log.Fatalf("Error updating device: %v", err)
						}

						updated, err := c.DecodeDeviceSource(res)
						if err != nil {
							log.Fatalf("Error adding device: %v", err)
						}

						log.Printf("Updated: %+v", updated)
					}

				}

				schema := fk.UpdateDeviceSourceSchemaPayload{
					Active:     true,
					Key:        streamName,
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
