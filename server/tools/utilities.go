package utilities

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	goaclient "github.com/goadesign/goa/client"

	fk "github.com/fieldkit/cloud/server/api/client"
)

func UpdateLocation(ctx context.Context, c *fk.Client, device *fk.DeviceSource, longitude, latitude float64) error {
	updateLocationPayload := fk.UpdateDeviceSourceLocationPayload{
		Longitude: longitude,
		Latitude:  latitude,
	}
	res, err := c.UpdateLocationDevice(ctx, fk.UpdateLocationDevicePath(device.ID), &updateLocationPayload)
	if err != nil {
		return fmt.Errorf("Error updating device location: %v", err)
	}

	updated, err := c.DecodeDeviceSource(res)
	if err != nil {
		return fmt.Errorf("Error adding device location: %v", err)
	}

	log.Printf("Updated: %+v", updated)

	return nil
}

func UpdateFirmware(ctx context.Context, c *fk.Client, deviceID int, firmwareID int) error {
	updatePayload := fk.UpdateDeviceFirmwarePayload{
		DeviceID:   deviceID,
		FirmwareID: firmwareID,
	}
	res, err := c.UpdateFirmware(ctx, fk.UpdateFirmwarePath(), &updatePayload)
	if err != nil {
		return fmt.Errorf("Error updating device firmware: %v", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Error updating device firmware: %v", res.Status)
	}

	log.Printf("Updated!")

	return nil
}

func FindExpedition(ctx context.Context, c *fk.Client, projectSlug string) (exp *fk.Expedition, err error) {
	res, err := c.ListProject(ctx, fk.ListProjectPath())
	if err != nil {
		return nil, fmt.Errorf("Error listing projects: %v", err)
	}
	projects, err := c.DecodeProjects(res)
	if err != nil {
		return nil, fmt.Errorf("Error listing projects: %v", err)
	}

	for _, project := range projects.Projects {
		log.Printf("Project: %+v\n", *project)

		if project.Slug == projectSlug {
			res, err = c.ListExpedition(ctx, fk.ListExpeditionPath(project.Slug))
			if err != nil {
				return nil, fmt.Errorf("%v", err)
			}
			expeditions, err := c.DecodeExpeditions(res)
			if err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			for _, exp := range expeditions.Expeditions {
				return exp, nil
			}
		}
	}

	return nil, fmt.Errorf("Unable to find expedition")
}

func FindExistingDevice(ctx context.Context, c *fk.Client, projectSlug, deviceKey string) (d *fk.DeviceSource, err error) {
	exp, err := FindExpedition(ctx, c, projectSlug)
	if err != nil {
		return nil, err
	}

	log.Printf("Expedition: %+v\n", *exp)

	res, err := c.ListDevice(ctx, fk.ListDevicePath(projectSlug, exp.Slug))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	devices, err := c.DecodeDeviceSources(res)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	for _, device := range devices.DeviceSources {
		if device.Key == deviceKey {
			log.Printf("Device: %+v\n", *device)
			return device, nil
		}
	}

	return nil, fmt.Errorf("Unable to find device")
}

func CreateWebDevice(ctx context.Context, c *fk.Client, projectSlug, deviceName, deviceId, streamName string) (d *fk.DeviceSource, err error) {
	theDevice, err := FindExistingDevice(ctx, c, projectSlug, deviceId)
	if err != nil {
		return nil, err
	}

	if theDevice == nil {
		log.Printf("Creating new Device %v", deviceName)
		exp, err := FindExpedition(ctx, c, projectSlug)
		if err != nil {
			return nil, err
		}

		addPayload := fk.AddDeviceSourcePayload{
			Name: deviceName,
			Key:  deviceId,
		}
		res, err := c.AddDevice(ctx, fk.AddDevicePath(exp.ID), &addPayload)
		if err != nil {
			return nil, fmt.Errorf("Error adding device: %v", err)
		}

		if res.StatusCode != 200 {
			bytes, _ := ioutil.ReadAll(res.Body)
			return nil, fmt.Errorf("Error adding device: %v", string(bytes))
		} else {
			added, err := c.DecodeDeviceSource(res)
			if err != nil {
				return nil, fmt.Errorf("Error adding device: %v", err)
			}

			theDevice = added
		}
	} else {
		if deviceId != "" && theDevice.Key != deviceId {
			log.Printf("Device id/key is different (%s != %s) fixing...", deviceId, theDevice.Key)

			updatePayload := fk.UpdateDeviceSourcePayload{
				Name: deviceName,
				Key:  deviceId,
			}
			res, err := c.UpdateDevice(ctx, fk.UpdateDevicePath(theDevice.ID), &updatePayload)
			if err != nil {
				return nil, fmt.Errorf("Error updating device: %v", err)
			}

			updated, err := c.DecodeDeviceSource(res)
			if err != nil {
				return nil, fmt.Errorf("Error adding device: %v", err)
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
		return nil, err
	}

	return theDevice, nil
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
