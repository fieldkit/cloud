package api

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"image"
	"image/jpeg"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type ResizedImage struct {
	Size        int64
	ContentType string
	Body        io.Reader
	Data        []byte
}

func resizeLoadedMedia(ctx context.Context, lm *repositories.LoadedMedia, newWidth, newHeight uint) (resized *ResizedImage, err error) {
	original, _, err := image.Decode(lm.Reader)
	if err != nil {
		return nil, err
	}

	resizer := nfnt.NewDefaultResizer()
	smaller := resizer.Resize(original, newWidth, newHeight)

	var buff bytes.Buffer

	jpeg.Encode(&buff, smaller, &jpeg.Options{
		Quality: 90,
	})

	resized = &ResizedImage{
		ContentType: "image/jpg",
		Size:        int64(buff.Len()),
		Body:        bytes.NewReader(buff.Bytes()),
		Data:        buff.Bytes(),
	}

	return
}

func smartCrop(original image.Image, cropX, cropY uint) (i image.Image, err error) {
	resizer := nfnt.NewDefaultResizer()
	analyzer := smartcrop.NewAnalyzer(resizer)
	topCrop, _ := analyzer.FindBestCrop(original, int(cropX), int(cropY))
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	cropped := original.(SubImager).SubImage(topCrop)
	thumb := resizer.Resize(cropped, cropX, cropY)
	return thumb, nil
}

func makeSimpleAssetURL(url string) string {
	return fmt.Sprintf("%s", url)
}

func makeAssetURL(url string, actual *string) *string {
	if actual == nil {
		return nil
	}
	hash := quickHash(*actual)
	final := fmt.Sprintf("%s?%s", url, hash)
	return &final
}

func quickHash(value string) string {
	h := sha1.New()
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

const ISO8601 = "2006-01-02T15:04:05-0700"
const RFC2822 = "Mon Jan 02 15:04:05 -0700 2006"

func tryParseDate(s *string) (time.Time, error) {
	if s == nil {
		return time.Time{}, fmt.Errorf("nil time string")
	}
	for _, layout := range []string{ISO8601, RFC2822, time.RFC3339} {
		if t, err := time.Parse(layout, *s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("parsing failed")
}
