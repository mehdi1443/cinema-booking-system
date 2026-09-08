package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp"
)

type Format int

const (
	FormatJPEG Format = iota + 1
)

func DecodeImage(r io.Reader) (image.Image, error) {
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
	case FormatJPEG:
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
			return nil, fmt.Errorf("failed to encode jpeg: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format type")
	}

	return &buf, nil
}

func GetContentType(format Format) string {
	switch format {
	case FormatJPEG:
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}
