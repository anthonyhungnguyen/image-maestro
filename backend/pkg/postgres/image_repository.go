package postgres

import (
	"encoding/json"

	"github.com/anthonyhungnguyen/image-maestro/model"
	"github.com/anthonyhungnguyen/image-maestro/pkg/config"
	"github.com/rs/zerolog/log"
)

func SaveImage(image_response model.ImageResponse) {
	db := config.GetConnection()
	defer db.Close()
	// Insert data
	sqlStatement := `
		INSERT INTO image (id, url, status, content_type, checksum, byte_size, width, height, exif)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	exifJson, err := json.Marshal(image_response.Exif)
	if err != nil {
		log.Fatal().Err(err).Msg("Error marshalling exif")
	}
	err = db.QueryRow(sqlStatement, image_response.Id, image_response.Url, image_response.Status, image_response.ContentType, image_response.Checksum, image_response.ByteSize, image_response.Width, image_response.Height, exifJson).Scan(&image_response.Id)
	if err != nil {
		log.Fatal().Err(err).Msg("Error inserting row")
	}

}
