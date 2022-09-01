package api

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"strings"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/files"
)

const (
	MimeTypeJpeg = "image/jpeg"
	MimeTypeGif  = "image/gif"
	MimeTypePng  = "image/png"
)

type PhotoCache struct {
	mr *repositories.MediaRepository
}

type ExternalMedia struct {
	URL         string
	ContentType string
}

type PhotoResizeSettings struct {
	Size int32
}

type PhotoCropSettings struct {
	X uint
	Y uint
}

type PhotoFromCache struct {
	Bytes       []byte
	Size        uint
	ContentType string
	Etag        string
	EtagMatch   bool
}

func NewPhotoCache(originals files.FileArchive) *PhotoCache {
	return &PhotoCache{
		mr: repositories.NewMediaRepository(originals),
	}
}

func (pc *PhotoCache) Load(ctx context.Context, media *ExternalMedia, resize *PhotoResizeSettings, crop *PhotoCropSettings, ifNoneMatch *string) (*PhotoFromCache, error) {
	etag := quickHash(media.URL)
	if resize != nil {
		etag += fmt.Sprintf(":%d", resize.Size)
	}

	if ifNoneMatch != nil {
		if *ifNoneMatch == fmt.Sprintf(`"%s"`, etag) {
			return &PhotoFromCache{
				Etag:        etag,
				Bytes:       []byte{},
				ContentType: media.ContentType,
				Size:        0,
				EtagMatch:   true,
			}, nil
		}
	}

	lm, err := pc.mr.LoadByURL(ctx, media.URL)
	if err != nil {
		return nil, err
	}

	haveGif := strings.Contains(media.ContentType, MimeTypeGif)

	if haveGif {
		buffer := make([]byte, lm.Size)
		_, err := io.ReadFull(lm.Reader, buffer)
		if err != nil {
			return nil, err
		}

		return &PhotoFromCache{
			Etag:        etag,
			Bytes:       buffer,
			ContentType: media.ContentType,
			Size:        uint(lm.Size),
		}, nil
	}

	if resize != nil && media.ContentType == MimeTypeJpeg || media.ContentType == MimeTypePng {
		resized, err := resizeLoadedMedia(ctx, lm, uint(resize.Size), 0)
		if err != nil {
			return nil, err
		}

		return &PhotoFromCache{
			Etag:        etag,
			Bytes:       resized.Data,
			ContentType: resized.ContentType,
			Size:        uint(resized.Size),
		}, nil
	} else if crop != nil {
		original, _, err := image.Decode(lm.Reader)
		if err != nil {
			return nil, nil
		}

		cropped, err := smartCrop(original, crop.X, crop.Y)
		if err != nil {
			return nil, err
		}

		options := jpeg.Options{
			Quality: 80,
		}

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, cropped, &options); err != nil {
			return nil, err
		}

		data := buf.Bytes()

		return &PhotoFromCache{
			Etag:        etag,
			Bytes:       data,
			ContentType: MimeTypeJpeg,
			Size:        uint(len(data)),
		}, nil
	} else {
		data, err := ioutil.ReadAll(lm.Reader)
		if err != nil {
			return nil, err
		}

		return &PhotoFromCache{
			Etag:        etag,
			Bytes:       data,
			ContentType: media.ContentType,
			Size:        uint(len(data)),
		}, nil
	}
}
