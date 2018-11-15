package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	fk "github.com/fieldkit/cloud/server/api/client"
	fktesting "github.com/fieldkit/cloud/server/tools"
)

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
	if len(fc.Firmwares) == 0 {
		return
	}

	sort.Sort(BySort(fc.Firmwares))

	w := tabwriter.NewWriter(os.Stdout, 8, 0, 4, ' ', tabwriter.Debug)
	fmt.Fprintln(w, fmt.Sprintf("ID\t Module\t Profile\t Time\t Tag\t Path"))

	for _, f := range fc.Firmwares {
		time := f.Time.Format("2006/01/02 15:04:05")

		if f.IsLocal() {
			fmt.Fprintln(w, fmt.Sprintf("%d\t %s\t %s\t %s\t %s\t %s", f.ID, f.Module, f.Profile, time, f.Etag, f.Path))
		} else {
			fmt.Fprintln(w, fmt.Sprintf("%d\t %s\t %s\t %s\t %s", f.ID, f.Module, f.Profile, time, f.Etag))
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

func (fm *FirmwareManager) ListLocal(ctx context.Context) (fc *FirmwareCollection, err error) {
	files, err := ioutil.ReadDir(fm.Directory)
	if err != nil {
		return nil, err
	}

	fc = NewFirmwareCollection()

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".bin.json") {
			path := path.Join(fm.Directory, file.Name())

			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}

			firmware := &Firmware{}
			err = json.Unmarshal(bytes, &firmware)
			if err != nil {
				return nil, err
			}

			fc.Add(firmware)
		}
	}

	return
}

func (fm *FirmwareManager) ListRemote(ctx context.Context, host, scheme string) (fc *FirmwareCollection, err error) {
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

func (fm *FirmwareManager) Find(ctx context.Context, id int) (f *Firmware, err error) {
	local, err := fm.ListLocal(ctx)
	if err != nil {
		return nil, err
	}

	for _, f := range local.Firmwares {
		if f.ID == id {
			return f, nil
		}
	}

	return nil, fmt.Errorf("No such firmware with ID: %v", id)
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
