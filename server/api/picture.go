package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"

	"github.com/goadesign/goa"
	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/O-C-R/fieldkit/server/api/app"
)

// InputController implements the picture resource.
type PictureController struct {
	*goa.Controller
}

func NewPictureController(service *goa.Service) *PictureController {
	return &PictureController{
		Controller: service.NewController("PictureController"),
	}
}

func (c *PictureController) GetID(ctx *app.GetIDPictureContext) error {
	hash := sha1.New()
	if err := binary.Write(hash, binary.LittleEndian, int64(ctx.UserID)); err != nil {
		return err
	}

	var r, g, b uint8
	hashBytesReader := bytes.NewReader(hash.Sum(nil))
	if err := binary.Read(hashBytesReader, binary.LittleEndian, &r); err != nil {
		return err
	}

	if err := binary.Read(hashBytesReader, binary.LittleEndian, &g); err != nil {
		return err
	}

	if err := binary.Read(hashBytesReader, binary.LittleEndian, &b); err != nil {
		return err
	}

	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.RGBA{r, g, b, 0xff})
	gc.MoveTo(0, 0)
	gc.LineTo(256, 0)
	gc.LineTo(256, 256)
	gc.LineTo(0, 256)
	gc.Close()
	gc.Fill()

	okResponse := bytes.NewBuffer([]byte{})
	if err := png.Encode(okResponse, img); err != nil {
		return err
	}

	return ctx.OK(okResponse.Bytes())
}
