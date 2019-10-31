package api

import (
    "bufio"
	"fmt"
    "io"

	"github.com/goadesign/goa"

    "github.com/fieldkit/cloud/server/backend/repositories"
)

func SendLoadedMedia(responseData *goa.ResponseData, lm *repositories.LoadedMedia) error {
    writer := bufio.NewWriter(responseData)

    responseData.Header().Set("Content-Length", fmt.Sprintf("%d", lm.Size))

    _, err := io.Copy(writer, lm.Reader)
    if err != nil {
        return err
    }

    writer.Flush()

    return nil
}
