package routes

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // register jpeg
	"io"
	"mime"
	"net/http"

	"github.com/anthonyhungnguyen/image-maestro/model"
	"github.com/anthonyhungnguyen/image-maestro/pkg/firebase"
	"github.com/anthonyhungnguyen/image-maestro/pkg/postgres"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/segmentio/ksuid"
)

func ImageRoutes(r *gin.RouterGroup) {
	r.POST("/image", getImage)
	r.POST("/image/:image_id/annotate", annotateImage)
	r.GET("/image/thumbnail", extractThumbnail)
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

	imageResponse := model.ImageResponse{
		Id:          ksuid.New().String(),
		Url:         requestBody.Url,
		ContentType: extension,
		Exif:        exif,
		Status:      "success",
		Checksum:    generateChecksum(body),
		Width:       img.Bounds().Dx(),
		Height:      img.Bounds().Dy(),
		ByteSize:    int64(len(body)),
	}
	// Upload image to firebase
	if err := firebase.UploadImage(body, imageResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "cannot upload image, exception: " + err.Error(),
		})
		c.Abort()
		return
	}

	// Save to database
	postgres.SaveImage(imageResponse)

	c.JSON(http.StatusOK, gin.H{
		"message": imageResponse,
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

func getToken(c *gin.Context) (string, error) {
	token, exists := c.Get("token")
	if exists {
		return token.(string), nil
	} else {
		log.Fatal().Msg("token not found")
		return "", errors.New("token not found")
	}
}

// TODO: implement this
func annotateImage(c *gin.Context) {
	imageId := c.Param("image_id")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("annotated image %s", imageId),
	})
}

func validateSig(c *gin.Context, imageId string, sig string) bool {
	token, exists := getToken(c)
	if exists != nil {
		return false
	}

	secret := []byte(token)
	message := []byte(imageId)

	// // Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(crypto.SHA256.New, secret)

	// Write Data to it
	h.Write(message)

	// Get result and encode as hexadecimal string
	signature := h.Sum(nil)
	hexSignature := hex.EncodeToString(signature)

	return hexSignature == sig
}

// Thumbnail
func extractThumbnail(c *gin.Context) {
	imageId := c.Query("image_id")
	// height := c.Query("height")
	// width := c.Query("width")
	sig := c.Query("sig")

	// Validate sig
	if !validateSig(c, imageId, sig) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
	}
	vips.Startup(nil)
	defer vips.Shutdown()

	// Generate thumbnail

	// imageId := c.Param("image_id")
	// log.Print(imageId)
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": fmt.Sprintf("extracted thumbnail for image %s", imageId),
	// })
}
