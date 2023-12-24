package firebase

import (
	"context"

	"github.com/anthonyhungnguyen/image-maestro/model"
	"github.com/anthonyhungnguyen/image-maestro/pkg/config"
	"github.com/rs/zerolog/log"
)

func UploadImage(imageBytes []byte, fileName model.ImageResponse) error {
	ctx := context.Background()

	// Get bucket
	bucket := config.GetBucket("image-maesto.appspot.com")

	obj := bucket.Object(fileName.Id)

	// Write file to bucket
	wc := obj.NewWriter(ctx)

	wc.ContentType = "image/" + fileName.ContentType[1:]
	log.Info().Msg(wc.ContentType)

	// Write content
	if _, err := wc.Write(imageBytes); err != nil {
		return err
	}

	// Close writer
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}
