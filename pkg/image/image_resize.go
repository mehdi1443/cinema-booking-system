package image

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/chai2010/webp"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/nfnt/resize"
	"github.com/pixiv/go-libjpeg/jpeg"
)

type Format int

const (
	FormatWebP Format = iota + 1
	FormatJPEG
)

func DecodeImage(r io.Reader) (image.Image, error) {
	image.RegisterFormat("webp", "RIFF????WEBP", webp.Decode, webp.DecodeConfig)

	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrInvalidFileType, err)
	}
	return img, nil
}

func Resize(width uint, height uint, img image.Image) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos2)
}

func EncodeToBuffer(img image.Image, format Format) (io.Reader, error) {
	var buf bytes.Buffer

	switch format {
	case FormatWebP:
		opts := &webp.Options{Lossless: false, Quality: 85}
		if err := webp.Encode(&buf, img, opts); err != nil {
			return nil, fmt.Errorf("failed to encode webp: %w", err)
		}
	case FormatJPEG:
		opts := &jpeg.EncoderOptions{Quality: 85, ProgressiveMode: true}
		if err := jpeg.Encode(&buf, img, opts); err != nil {
			return nil, fmt.Errorf("failed to encode jpeg: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format type")
	}

	return &buf, nil
}

func GetContentType(format Format) string {
	switch format {
	case FormatWebP:
		return "image/webp"
	case FormatJPEG:
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}
