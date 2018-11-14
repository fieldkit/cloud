package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
)

var (
	active    = 0
	failed    = 0
	succeeded = 0
)

func DownloadAll(fc *FirmwareCollection) error {
	client := grab.NewClient()
	reqs := make([]*grab.Request, len(fc.Firmwares))
	for i, f := range fc.Firmwares {
		req, err := grab.NewRequest(".", f.URL)
		if err != nil {
			return err
		}
		req.Filename = f.Path
		reqs[i] = req
	}

	fmt.Printf("Downloading %d files...\n", len(fc.Firmwares))
	respch := client.DoBatch(3, reqs...)

	responses := make([]*grab.Response, 0, len(fc.Firmwares))
	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case resp := <-respch:
			if resp != nil {
				responses = append(responses, resp)
			} else {
				// Channel is closed - all downloads are complete
				update(responses)
				break Loop
			}

		case <-t.C:
			update(responses)
		}
	}

	fmt.Printf(
		"Finished %d successful, %d failed, %d incomplete.\n",
		succeeded,
		failed,
		active)

	return nil
}

func update(responses []*grab.Response) {
	if active > 0 {
		fmt.Printf("\033[%dA\033[K", active)
	}

	for i, resp := range responses {
		if resp != nil && resp.IsComplete() {
			if resp.Err() != nil {
				failed++
				fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", resp.Request.URL(), resp.Err())
			} else {
				succeeded++
				fmt.Printf("Finished %s %d / %d bytes (%d%%)\n",
					resp.Filename,
					resp.BytesComplete(),
					resp.Size,
					int(100*resp.Progress()))
			}
			responses[i] = nil
		}
	}

	active = 0
	for _, resp := range responses {
		if resp != nil {
			fmt.Printf("Downloading %s %d / %d bytes (%d%%) - %.02fKBp/s ETA: %ds \033[K\n",
				resp.Filename,
				resp.BytesComplete(),
				resp.Size,
				int(100*resp.Progress()),
				resp.BytesPerSecond()/1024,
				int64(resp.ETA().Sub(time.Now()).Seconds()))
			active++
		}
	}
}
