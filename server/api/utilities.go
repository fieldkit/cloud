package api

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"image"
	_ "image/color"
	_ "image/draw"
	"image/jpeg"

	"github.com/goadesign/goa"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

func sendLoadedMedia(responseData *goa.ResponseData, lm *repositories.LoadedMedia) error {
	writer := bufio.NewWriter(responseData)

	responseData.Header().Set("Content-Length", fmt.Sprintf("%d", lm.Size))

	_, err := io.Copy(writer, lm.Reader)
	if err != nil {
		return err
	}

	writer.Flush()

	return nil
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

func sendImage(responseData *goa.ResponseData, image image.Image) error {
	options := jpeg.Options{
		Quality: 80,
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, &options)
	if err != nil {
		return err
	}

	return sendData(responseData, "image/jpeg", buf.Bytes())
}

func logErrorAndSendData(ctx context.Context, responseData *goa.ResponseData, cause error, contentType string, data []byte) error {
	log := Logger(ctx).Sugar()

	log.Errorw("error", "error", cause)

	return sendData(responseData, contentType, data)
}

func sendData(responseData *goa.ResponseData, contentType string, data []byte) error {
	responseData.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	responseData.Header().Set("Content-Type", contentType)
	responseData.Write(data)
	return nil
}

func makePhotoURL(url string, actual *string) *string {
	if actual == nil {
		return nil
	}
	h := sha1.New()
	h.Write([]byte(*actual))
	hash := hex.EncodeToString(h.Sum(nil))
	final := fmt.Sprintf("%s?%s", url, hash)
	return &final
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
