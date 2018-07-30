package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	fk "github.com/fieldkit/cloud/server/api/client"
	fktesting "github.com/fieldkit/cloud/server/tools"
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

	Latitude         float64
	Longitude        float64
	LocationPrimeZip string

	FirmwareURL  string
	FirmwareETag string

	Module string

	FirmwareDirectory string
	FirmwareFile      string
	FirmwareMeta      string
}

type Metadata struct {
	Module  string
	ETag    string
	Profile string
	Map     map[string]*string
}

func getModuleFromJobName(name string) (string, error) {
	return strings.Replace(name, "/", "-", 1), nil
}

func getProfileFromFile(module, path string) (string, error) {
	file := filepath.Base(path)
	re := regexp.MustCompile(fmt.Sprintf("%s-(.+).bin", module))
	m := re.FindAllStringSubmatch(file, -1)
	if len(m) == 0 {
		return "", fmt.Errorf("Malformed file name %s (%s), no profile for %s", path, file, module)
	}
	return m[0][1], nil
}

func getMetaFromEnvironment(moduleOverride string, file string) (metadata *Metadata, err error) {
	jobName := os.Getenv("JOB_NAME")
	if jobName == "" {
		return nil, fmt.Errorf("ENV[JOB_NAME] missing.")
	}

	module := moduleOverride
	if moduleOverride == "" {
		module, err = getModuleFromJobName(jobName)
		if err != nil {
			return nil, fmt.Errorf("Error getting module from job name: %v", err)
		}
		log.Printf("Found module name: '%s'", module)
	} else {
		log.Printf("Using module override: '%s'", module)
	}

	buildTime := os.Getenv("BUILD_TIMESTAMP")
	layout := "20060102_150405"
	_, err = time.Parse(layout, buildTime)
	if err != nil {
		return nil, fmt.Errorf("ENV[BUILD_TIMESTAMP] parse failed.")
	}

	profile, err := getProfileFromFile(module, file)
	if err != nil {
		return nil, err
	}
	etag := fmt.Sprintf("%s_%s_%s", module, profile, buildTime)

	metadata = &Metadata{
		ETag:    etag,
		Module:  module,
		Profile: profile,
		Map:     make(map[string]*string),
	}

	metadata.Map["Build-Id"] = aws.String(os.Getenv("BUILD_ID"))
	metadata.Map["Build-Number"] = aws.String(os.Getenv("BUILD_NUMBER"))
	metadata.Map["Build-Tag"] = aws.String(os.Getenv("BUILD_TAG"))
	metadata.Map["Build-Time"] = aws.String(buildTime)
	metadata.Map["Build-Commit"] = aws.String(os.Getenv("GIT_COMMIT"))
	metadata.Map["Build-Branch"] = aws.String(os.Getenv("GIT_BRANCH"))
	metadata.Map["Build-Job-Base"] = aws.String(os.Getenv("JOB_BASE_NAME"))
	metadata.Map["Build-Job-Name"] = aws.String(jobName)
	metadata.Map["Build-Module"] = aws.String(module)

	return
}

func uploadAllFirmware(ctx context.Context, c *fk.Client, moduleOverride, directory string) error {
	files, err := filepath.Glob(directory + "/*.bin")
	if err != nil {
		return err
	}
	for _, filename := range files {
		log.Printf("Processing %s...", filename)
		err := uploadFirmware(ctx, c, moduleOverride, filename)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
	return nil
}

func createAwsSession() (s *session.Session, err error) {
	configs := []aws.Config{
		aws.Config{
			Region: aws.String("us-east-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		aws.Config{
			Region:                        aws.String("us-east-1"),
			Credentials:                   credentials.NewEnvCredentials(),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}

	for _, config := range configs {
		sessionOptions := session.Options{
			Profile: "fieldkit",
			Config:  config,
		}
		session, err := session.NewSessionWithOptions(sessionOptions)
		if err == nil {
			return session, nil
		}
	}

	return nil, fmt.Errorf("Error creating AWS session: %v", err)
}

func uploadFirmware(ctx context.Context, c *fk.Client, moduleOverride, filename string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("Error creating UUID: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", filename)
	}

	defer file.Close()

	metadata, err := getMetaFromEnvironment(moduleOverride, filename)
	if err != nil {
		return err
	}

	session, err := createAwsSession()
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(session)

	log.Printf("Uploading %s...", filename)

	r, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         nil,
		ContentType: aws.String("application/octet-stream"),
		Bucket:      aws.String("conservify-firmware"),
		Key:         aws.String(id.String()),
		Body:        file,
		Metadata:    metadata.Map,
		Tagging:     nil,
	})
	if err != nil {
		return fmt.Errorf("Error uploading firmware: %v", err)
	}

	log.Printf("Uploaded %s", r.Location)

	log.Printf("Creating database entry...")

	jsonData, err := json.Marshal(metadata.Map)
	if err != nil {
		return fmt.Errorf("Error serializing metadata: %v", err)
	}

	addFirmwarePayload := fk.AddFirmwarePayload{
		Etag:   metadata.ETag,
		Module: metadata.Module,
		URL:    r.Location,
		Meta:   string(jsonData),
	}
	_, err = c.AddFirmware(ctx, fk.AddFirmwarePath(), &addFirmwarePayload)
	if err != nil {
		return fmt.Errorf("Error adding firmware: %v", err)
	}

	log.Printf("Done!")

	return nil
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

	flag.Float64Var(&o.Longitude, "longitude", 0, "longitude")
	flag.Float64Var(&o.Latitude, "latitude", 0, "latitude")
	flag.StringVar(&o.LocationPrimeZip, "location-prime-zip", "", "zipcode to prime the location with")

	flag.StringVar(&o.FirmwareURL, "firmware-url", "", "firmware url")
	flag.StringVar(&o.FirmwareETag, "firmware-etag", "", "firmware etag")

	flag.StringVar(&o.Module, "module", "", "override module")

	flag.StringVar(&o.FirmwareDirectory, "firmware-directory", "", "firmware directory")
	flag.StringVar(&o.FirmwareFile, "firmware-file", "", "firmware file")
	flag.StringVar(&o.FirmwareMeta, "firmware-meta", "", "firmware meta")

	flag.Parse()

	c, err := fktesting.CreateAndAuthenticate(ctx, o.Host, o.Scheme, o.Username, o.Password)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("Authenticated as %s (%s)", o.Username, o.Host)

	if o.FirmwareFile != "" {
		err := uploadFirmware(ctx, c, o.Module, o.FirmwareFile)
		if err != nil {
			log.Fatalf("Error adding firmware: %v", err)
		}
	}

	if o.FirmwareDirectory != "" {
		err := uploadAllFirmware(ctx, c, o.Module, o.FirmwareDirectory)
		if err != nil {
			log.Fatalf("Error adding firmware: %v", err)
		}
	}

	if o.DeviceName != "" {
		device, err := fktesting.CreateWebDevice(ctx, c, o.Project, o.DeviceName, o.DeviceId, "")
		if err != nil {
			log.Fatalf("Error creating web device: %v", err)
		}

		log.Printf("Associated %v", device)

		if o.Latitude != 0 && o.Longitude != 0 {
			log.Printf("Setting location %v,%v", o.Longitude, o.Latitude)
			err := fktesting.UpdateLocation(ctx, c, device, o.Longitude, o.Latitude)
			if err != nil {
				log.Fatalf("Error updating location: %v", err)
			}
		}

		if o.FirmwareURL != "" && o.FirmwareETag != "" {
			log.Printf("Updating firmware %v", o.FirmwareURL)
			err := fktesting.UpdateFirmware(ctx, c, device, o.FirmwareURL, o.FirmwareETag)
			if err != nil {
				log.Fatalf("Error updating firmware: %v", err)
			}
		}
	}
}
