package routes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	_ "image/jpeg" // register jpeg
	"io"
	"mime"
	"net/http"

	"github.com/anthonyhungnguyen/image-maestro/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/segmentio/ksuid"
)

func ImageRoutes(r *gin.RouterGroup) {
	r.POST("/image", getImage)
}

func getImage(c *gin.Context) {
	log.Print(c.Request.Body)

	var requestBody model.ImageRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		c.Abort()
		return
	}

	imageResp, err := http.Get(requestBody.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot pull image, exception: " + err.Error(),
		})
		c.Abort()
		return
	}

	defer imageResp.Body.Close()

	// Read body to process later
	body, err := io.ReadAll(imageResp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot read image, exception: " + err.Error(),
		})
		c.Abort()
		return
	}

	extension := detectMIMEType(body)

	if extension != ".jpe" && extension != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "only support jpeg",
		})
		c.Abort()
		return
	}

	bytesReader := bytes.NewReader(body)

	// Decode image
	img, _, err := image.Decode(bytesReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot decode image, exception: " + err.Error(),
		})
		c.Abort()
		return
	}

	bytesReader = bytes.NewReader(body)

	// Extract exif
	exif, err := exif.Decode(bytesReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "cannot extract exif, exception: " + err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": model.ImageResponse{
			Id:          ksuid.New().String(),
			Url:         requestBody.Url,
			ContentType: extension,
			Exif:        exif,
			Status:      "success",
			Checksum:    generateChecksum(body),
			Width:       img.Bounds().Dx(),
			Height:      img.Bounds().Dy(),
			ByteSize:    int64(len(body)),
		},
	})
}

func detectMIMEType(content []byte) string {
	// Read the first 512 bytes to pass to DetectContentType
	buffer := content[:512]

	// Detect the content type of the file
	fileType := http.DetectContentType(buffer)

	// Infer extension
	extensions, err := mime.ExtensionsByType(fileType)

	if err != nil {
		log.Fatal().Err(err).Msg("error getting extension")
	}

	// Return the first extension
	return extensions[0]
}

func generateChecksum(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
