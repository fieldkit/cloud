package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	tooling "github.com/Conservify/tooling"
)

type options struct {
	CacheDirectory string
	Scheme         string
	Host           string

	Download bool
	Offline  bool
	All      bool
	Module   string
	Profile  string

	ListPorts  bool
	Port       string
	FirmwareID int
}

func main() {
	ctx := context.TODO()

	o := &options{}

	flag.StringVar(&o.CacheDirectory, "cache-directory", getDefaultCacheDirectory(), "cache directory")

	flag.StringVar(&o.Scheme, "scheme", "https", "fk instance scheme")
	flag.StringVar(&o.Host, "host", "api.fkdev.org", "fk instance hostname")

	flag.StringVar(&o.Module, "module", "", "filter by module")
	flag.StringVar(&o.Profile, "profile", "", "filter by profile")
	flag.BoolVar(&o.All, "all", false, "disable filtering by latest firmware")

	flag.BoolVar(&o.Offline, "offline", false, "work in offline mode")

	flag.BoolVar(&o.ListPorts, "ports", false, "list ports")

	flag.BoolVar(&o.Download, "download", false, "download firmware and cache locally")

	flag.IntVar(&o.FirmwareID, "flash", 0, "flash firmware to a connected device")
	flag.StringVar(&o.Port, "port", "", "port to focus on when flashing")

	flag.Parse()

	fm, err := NewFirmwareManager(ctx, o.CacheDirectory, o.Host, o.Scheme)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	if o.ListPorts {
		pd := tooling.NewPortDiscoveror()
		for _, p := range pd.List() {
			log.Printf("Port: %s", p.Name)
		}
	}

	if o.FirmwareID > 0 {
		f, err := fm.Find(ctx, o.FirmwareID, o.Offline)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		ae := tooling.NewArduinoEnvironment()

		err = ae.Locate("")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		log.Printf("Preparing to flash new firmware. Please quickly double-press RST button on the board.")

		err = tooling.Upload(&tooling.UploadOptions{
			Arduino:     ae,
			SkipTouch:   false,
			Board:       "adafruit_feather_m0",
			Binary:      f.Path,
			Port:        "",
			FlashOffset: 32768,
			Verbose:     false,
			Verify:      true,
			Quietly:     false,
		})
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		log.Printf("Success!")
	} else {
		fc, err := getFirmwareCollection(ctx, fm, o)
		if err != nil {
			log.Fatalf("Error: %v", err)
			return
		}

		if !o.All {
			fc = fc.FilterLatest()
		}

		if o.Module != "" {
			fc = fc.FilterModule(o.Module)
		}

		if o.Profile != "" {
			fc = fc.FilterProfile(o.Profile)
		}

		if o.Download {
			err := fm.Download(fc.NotLocal())
			if err != nil {
				log.Fatalf("Error: %v", err)
				return
			}
		}

		if fc.Empty() {
			fmt.Println("No firmware found matching that criteria.")
		} else {
			fc.Display()
		}
	}
}

func getDefaultCacheDirectory() string {
	home := os.Getenv("HOME")
	if home == "" {
		panic("Unable to find HOME")
	}
	return path.Join(home, ".fk/firmware")
}

func getFirmwareCollection(ctx context.Context, fm *FirmwareManager, o *options) (*FirmwareCollection, error) {
	fc, err := fm.ListLocal(ctx)
	if err != nil {
		return nil, err
	}

	if !o.Offline {
		log.Printf("Querying %s...", o.Host)

		remote, err := fm.ListRemote(ctx)
		if err != nil {
			return nil, err
		}

		fc = fc.Merge(remote)
	}

	return fc, nil
}
