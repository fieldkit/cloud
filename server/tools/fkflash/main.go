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
	Scheme         string
	Host           string
	ListRemote     bool
	ListLocal      bool
	ListPorts      bool
	Download       bool
	Offline        bool
	All            bool
	Module         string
	CacheDirectory string
	FirmwareID     int
	Port           string
}

func main() {
	ctx := context.TODO()

	o := options{}

	flag.StringVar(&o.CacheDirectory, "cache-directory", getDefaultCacheDirectory(), "cache directory")

	flag.StringVar(&o.Scheme, "scheme", "https", "fk instance scheme")
	flag.StringVar(&o.Host, "host", "api.fkdev.org", "fk instance hostname")

	flag.StringVar(&o.Module, "module", "", "module")
	flag.BoolVar(&o.All, "all", false, "disable filtering by latest firmware")

	flag.BoolVar(&o.Offline, "offline", false, "work in offline mode")

	flag.BoolVar(&o.ListRemote, "remote", false, "list firmware")
	flag.BoolVar(&o.ListLocal, "local", false, "list firmware")
	flag.BoolVar(&o.ListPorts, "ports", false, "list ports")

	flag.BoolVar(&o.Download, "download", false, "download firmware and cache locally")

	flag.IntVar(&o.FirmwareID, "flash", 0, "flash firmware to a connected device")
	flag.StringVar(&o.Port, "port", "", "port to focus on when flashing")

	flag.Parse()

	firmware, err := NewFirmwareManager(ctx, o.CacheDirectory)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	if !o.Offline {
	}

	if o.ListRemote {
		fc, err := firmware.ListRemote(ctx, o.Host, o.Scheme)
		if err != nil {
			log.Fatalf("Error: %v", err)
			return
		}

		if !o.All {
			fc = fc.Latest()
		}

		if o.Download {
			err := firmware.Download(fc.NotLocal())
			if err != nil {
				log.Fatalf("Error: %v", err)
				return
			}
		}

		fmt.Println("REMOTE:")

		fc.Display()
	}

	if o.Download {
		fc, err := firmware.ListRemote(ctx, o.Host, o.Scheme)
		if err != nil {
			log.Fatalf("Error: %v", err)
			return
		}

		if !o.All {
			fc = fc.Latest()
		}

		err = firmware.Download(fc.NotLocal())
		if err != nil {
			log.Fatalf("Error: %v", err)
			return
		}
	}

	if o.ListLocal {
		fc, err := firmware.ListLocal(ctx)
		if err != nil {
			log.Fatalf("Error: %v", err)
			return
		}

		if fc.Empty() {
			log.Printf("No local firmware. Try running with --download if you're online.")
			return
		}

		if !o.All {
			fc = fc.Latest()
		}

		fmt.Println("LOCAL:")

		fc.Display()
	}

	if o.ListPorts {
		pd := tooling.NewPortDiscoveror()
		for _, p := range pd.List() {
			fmt.Println("Port:", p.Name)
		}
	}

	if o.FirmwareID > 0 {
		f, err := firmware.Find(ctx, o.FirmwareID)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		ae := tooling.NewArduinoEnvironment()

		err = ae.Locate("")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		tooling.Upload(&tooling.UploadOptions{
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
	}
}

func getDefaultCacheDirectory() string {
	home := os.Getenv("HOME")
	if home == "" {
		panic("Unable to find HOME")
	}
	return path.Join(home, ".fk/firmware")
}
