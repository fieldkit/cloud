package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	fk "github.com/fieldkit/cloud/server/api/client"
	fktesting "github.com/fieldkit/cloud/server/tools"
)

import "text/tabwriter"

type options struct {
	Scheme         string
	Host           string
	List           bool
	Download       bool
	All            bool
	Module         string
	CacheDirectory string
}

func getDefaultCacheDirectory() string {
	return ".fk/firmware"
}

func main() {
	ctx := context.TODO()

	o := options{}

	flag.StringVar(&o.CacheDirectory, "cache-directory", getDefaultCacheDirectory(), "cache directory")

	flag.StringVar(&o.Scheme, "scheme", "https", "fk instance scheme")
	flag.StringVar(&o.Host, "host", "api.fkdev.org", "fk instance hostname")
	flag.StringVar(&o.Module, "module", "", "override module")

	flag.BoolVar(&o.List, "list", false, "list firmware")
	flag.BoolVar(&o.All, "all", false, "disable filtering by latest firmware")
	flag.BoolVar(&o.Download, "download", false, "download firmware and cache locally")

	flag.Parse()

	firmware, err := NewFirmwareManager(ctx, o.CacheDirectory)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	if o.List || o.Download {
		fc, err := firmware.List(ctx, o.Host, o.Scheme)
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

		fc.Display()
	}
}

type FirmwareManager struct {
	Directory string
}

func NewFirmwareManager(ctx context.Context, directory string) (fm *FirmwareManager, err error) {
	fm = &FirmwareManager{
		Directory: directory,
	}

	err = fm.Initialize(ctx)
	if err != nil {
		return nil, err
	}

	return
}

func (fm *FirmwareManager) Initialize(ctx context.Context) error {
	return mkdirAllIfMissing(fm.Directory)
}

func (fm *FirmwareManager) List(ctx context.Context, host, scheme string) (fc *FirmwareCollection, err error) {
	c, err := fktesting.NewClient(ctx, host, scheme)
	if err != nil {
		return nil, err
	}

	res, err := c.ListFirmware(ctx, fk.ListFirmwarePath(), nil, nil)
	if err != nil {
		return nil, err
	}
	firmwares, err := c.DecodeFirmwares(res)
	if err != nil {
		return nil, err
	}

	fc = NewFirmwareCollection()

	for _, f := range firmwares.Firmwares {
		fc.Add(&Firmware{
			ID:      f.ID,
			Module:  f.Module,
			Profile: f.Profile,
			Time:    f.Time,
			Etag:    f.Etag,
			URL:     f.URL,
			Path:    fmt.Sprintf("%s/%s.bin", fm.Directory, f.Etag),
		})
	}

	return fc, nil
}

func (fm *FirmwareManager) Download(fc *FirmwareCollection) error {
	if fc.Empty() {
		return nil
	}
	return DownloadAll(fc)
}

func mkdirAllIfMissing(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = os.MkdirAll(p, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

type FirmwareCollection struct {
	Firmwares []*Firmware
}

type Firmware struct {
	ID       int
	Module   string
	Profile  string
	Time     time.Time
	Etag     string
	URL      string
	Filename string
	Path     string
}

type BySort []*Firmware

func (a BySort) Len() int           { return len(a) }
func (a BySort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySort) Less(i, j int) bool { return a[i].Time.Unix() > a[j].Time.Unix() }

func NewFirmwareCollection() (fc *FirmwareCollection) {
	return &FirmwareCollection{
		Firmwares: make([]*Firmware, 0),
	}
}

func (fc *FirmwareCollection) NotLocal() *FirmwareCollection {
	newFc := NewFirmwareCollection()

	for _, f := range fc.Firmwares {
		if f.Path != "" {
			if !f.IsLocal() {
				newFc.Add(f)
			}
		}
	}

	return newFc
}

func (fc *FirmwareCollection) Latest() *FirmwareCollection {
	newFc := NewFirmwareCollection()

	seen := make(map[string]bool)

	sort.Sort(BySort(fc.Firmwares))

	for _, f := range fc.Firmwares {
		key := f.Module + "-" + f.Profile
		if _, ok := seen[key]; !ok {
			newFc.Firmwares = append(newFc.Firmwares, f)
			seen[key] = true
		}
	}

	return newFc
}

func (fc *FirmwareCollection) Empty() bool {
	return len(fc.Firmwares) == 0
}

func (fc *FirmwareCollection) Add(f *Firmware) {
	fc.Firmwares = append(fc.Firmwares, f)
}

func (fc *FirmwareCollection) Display() {
	w := tabwriter.NewWriter(os.Stdout, 8, 0, 4, ' ', tabwriter.Debug)
	fmt.Fprintln(w, fmt.Sprintf("ID\t Module\t Profile\t Time\t Tag"))

	for _, f := range fc.Firmwares {
		if f.IsLocal() {
			fmt.Fprintln(w, fmt.Sprintf("%d\t %s\t %s\t %s\t %s\t %s", f.ID, f.Module, f.Profile, f.Time, f.Etag, f.Path))
		} else {
			fmt.Fprintln(w, fmt.Sprintf("%d\t %s\t %s\t %s\t %s", f.ID, f.Module, f.Profile, f.Time, f.Etag))
		}
	}

	w.Flush()
}

func (f *Firmware) IsLocal() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}
