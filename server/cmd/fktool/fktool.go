package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type options struct {
	Scheme   string
	Host     string
	Email    string
	Password string

	Version string

	Module  string
	Profile string

	FirmwareFile string

	DryRun bool
}

type Metadata struct {
	ETag    string
	Module  string
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
		return "", fmt.Errorf("malformed file name %s (%s), no profile for %s", path, file, module)
	}
	return m[0][1], nil
}

func getMetaFromEnvironment(moduleOverride, profileOverride, version, file string) (metadata *Metadata, err error) {
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
		log.Printf("found module name: '%s'", module)
	} else {
		log.Printf("using module override: '%s'", module)
	}

	buildTime := os.Getenv("BUILD_TIMESTAMP")
	layout := "20060102_150405"
	_, err = time.Parse(layout, buildTime)
	if err != nil {
		return nil, fmt.Errorf("ENV[BUILD_TIMESTAMP] parse failed.")
	}

	profile := profileOverride
	if profile == "" {
		profile, err = getProfileFromFile(module, file)
		if err != nil {
			return nil, err
		}
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
	metadata.Map["Build-Profile"] = aws.String(profile)
	metadata.Map["Build-Version"] = aws.String(version)

	return
}

func createAwsSession() (s *session.Session, err error) {
	configs := []aws.Config{
		aws.Config{
			Region:                        aws.String("us-east-1"),
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
			Config: config,
		}
		session, err := session.NewSessionWithOptions(sessionOptions)
		if err == nil {
			return session, nil
		}
	}

	return nil, fmt.Errorf("error creating AWS session: %v", err)
}

func getFileHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	h := hex.EncodeToString(hasher.Sum(nil))

	return h, nil
}

func hasFile(session *session.Session, id string) (string, error) {
	bucket := "conservify-firmware"

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(id),
	}

	svc := s3.New(session)

	_, err := svc.HeadObject(hoi)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // TODO Make this s3.ErrCodeNoSuchKey eventually.
				return "", nil
			}
		}

		return "", err
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, id), nil
}

func uploadFirmware(ctx context.Context, fkc *FkClient, moduleOverride, profileOverride, version, filename string, dryRun bool) error {
	id, err := getFileHash(filename)
	if err != nil {
		return fmt.Errorf("error getting file hash: %v", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", filename)
	}

	defer file.Close()

	metadata, err := getMetaFromEnvironment(moduleOverride, profileOverride, version, filename)
	if err != nil {
		return err
	}

	session, err := createAwsSession()
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(session)

	url := ""

	existingUrl, err := hasFile(session, id)
	if err != nil {
		return err
	}

	if existingUrl != "" {
		log.Printf("object already uploaded")
		url = existingUrl
	} else {
		log.Printf("uploading %s...", filename)

		if !dryRun {
			r, err := uploader.Upload(&s3manager.UploadInput{
				ACL:         nil,
				ContentType: aws.String("application/octet-stream"),
				Bucket:      aws.String("conservify-firmware"),
				Key:         aws.String(id),
				Body:        file,
				Metadata:    metadata.Map,
				Tagging:     nil,
			})
			if err != nil {
				return fmt.Errorf("error uploading firmware: %v", err)
			}

			log.Printf("uploaded %s", r.Location)

			url = r.Location
		}
	}

	log.Printf("creating database entry for %s...", version)

	jsonData, err := json.Marshal(metadata.Map)
	if err != nil {
		return fmt.Errorf("error serializing metadata: %v", err)
	}

	addFirmwarePayload := AddFirmwarePayload{
		Etag:    metadata.ETag,
		Module:  metadata.Module,
		Profile: metadata.Profile,
		Version: version,
		URL:     url,
		Meta:    string(jsonData),
	}

	if !dryRun {
		err := fkc.AddFirmware(ctx, &addFirmwarePayload)
		if err != nil {
			return fmt.Errorf("error adding firmware: %v", err)
		}

		log.Printf("added")
	}

	log.Printf("done!")

	return nil
}

func main() {
	ctx := context.TODO()

	o := options{}

	flag.StringVar(&o.Scheme, "scheme", "http", "fk instance scheme")
	flag.StringVar(&o.Host, "host", "127.0.0.1:8080", "fk instance hostname")
	flag.StringVar(&o.Email, "email", "info@conservify.org", "email")
	flag.StringVar(&o.Password, "password", "asdfasdfasdf", "password")
	flag.StringVar(&o.Version, "version", "", "version")
	flag.StringVar(&o.Module, "module", "", "override module")
	flag.StringVar(&o.Profile, "profile", "", "override profile")
	flag.StringVar(&o.FirmwareFile, "firmware-file", "", "firmware file")
	flag.BoolVar(&o.DryRun, "dry", false, "dry run")

	flag.Parse()

	if o.Version == "" {
		log.Fatalf("version is required")
	}

	fkc := NewFkClient(o.Host, o.Scheme)

	err := fkc.Login(ctx, o.Email, o.Password)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("authenticated as %s (%s)", o.Email, o.Host)

	if o.FirmwareFile != "" {
		err := uploadFirmware(ctx, fkc, o.Module, o.Profile, o.Version, o.FirmwareFile, o.DryRun)
		if err != nil {
			log.Fatalf("error adding firmware: %v", err)
		}
	}
}

type FkClient struct {
	scheme string
	host   string
	http   *http.Client
	auth   string
}

func NewFkClient(host, scheme string) (fkc *FkClient) {
	return &FkClient{
		scheme: scheme,
		host:   host,
		http:   http.DefaultClient,
	}
}

func (fkc *FkClient) Login(ctx context.Context, email, password string) error {
	payload := &LoginPayload{
		Email:    email,
		Password: password,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s://%s/%s", fkc.scheme, fkc.host, "login")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	response, err := fkc.http.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid username or password")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fkc.auth = response.Header.Get("Authorization")

	_ = body

	return nil
}

func (fkc *FkClient) AddFirmware(ctx context.Context, payload *AddFirmwarePayload) error {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s://%s/%s", fkc.scheme, fkc.host, "firmware")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", fkc.auth)
	if err != nil {
		return err
	}

	response, err := fkc.http.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return fmt.Errorf("error adding firmware")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	_ = body

	return nil
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddFirmwarePayload struct {
	Etag    string `json"etag"`
	Module  string `json:"module"`
	Profile string `json:"profile"`
	Version string `json:"version"`
	URL     string `json:"url"`
	Meta    string `json:"meta"`
}
