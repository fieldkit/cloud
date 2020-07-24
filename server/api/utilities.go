package api

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"image"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
)

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

type Signer struct {
	key []byte
}

func NewSigner(key []byte) (s *Signer) {
	return &Signer{
		key: key,
	}
}

func (s *Signer) SignURL(url string) (string, error) {
	if true {
		return url, nil
	}
	now := time.Now()
	token := jwtgo.New(jwtgo.SigningMethodHS512)
	token.Claims = jwtgo.MapClaims{
		"exp": now.Add(time.Hour * 1).Unix(),
	}

	signed, err := token.SignedString(s.key)
	if err != nil {
		return "", err
	}

	return url + "?token=" + signed, nil
}

func (s *Signer) SignAndBustURL(url string, key *string) (*string, error) {
	if key == nil {
		return nil, nil
	}

	signed, err := s.SignURL(url)
	if err != nil {
		return nil, err
	}

	return &signed, nil
}
