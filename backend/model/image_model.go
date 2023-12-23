package model

import "github.com/rwcarlsen/goexif/exif"

type ImageRequest struct {
	Url  string `json:"url"`
	Exif bool   `json:"exif"`
}

type ImageResponse struct {
	Id          string     `json:"id"`
	Url         string     `json:"url"`
	Status      string     `json:"status"`
	ContentType string     `json:"content_type"`
	Checksum    string     `json:"checksum"`
	ByteSize    int64      `json:"byte_size"`
	Width       int        `json:"width"`
	Height      int        `json:"height"`
	Exif        *exif.Exif `json:"exif"`
}
