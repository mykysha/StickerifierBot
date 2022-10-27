// Package photoconverter makes photo manipulations.
package photoconverter

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
)

var ErrUnsupportedFormat = errors.New("unsupported photo format")

// ConvertToPNG converts a photo to PNG.
func ConvertToPNG(photo []byte) ([]byte, error) {
	/*
		contentType := http.DetectContentType(photo)

		switch {
		case contentType == "image/png":
			resized, err := resizePNG(photo, 512)
			if err != nil {
				return nil, fmt.Errorf("unable to resize png: %w", err)
			}

			return resized, nil
		case contentType == "image/jpeg":
			photo, err := convertJPEGToPNG(photo)
			if err != nil {
				return nil, fmt.Errorf("unable to convert jpeg to png: %w", err)
			}
	*/
	resized, err := resizePNG(photo, 512)
	if err != nil {
		return nil, fmt.Errorf("unable to resize png: %w", err)
	}

	return resized, nil
	/*
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, contentType)
		}*/
}

// convertJPEGToPNG converts a JPEG photo to PNG.
func convertJPEGToPNG(photo []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(photo))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, err)
		}

		return nil, fmt.Errorf("unable to decode jpeg: %w", err)
	}

	buf := new(bytes.Buffer)

	err = png.Encode(buf, img)
	if err != nil {
		return nil, fmt.Errorf("unable to encode png: %w", err)
	}

	return buf.Bytes(), nil
}

// resizePNG resizes a png photo so that the biggest side is equal to the given size.
func resizePNG(photo []byte, size int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(photo))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, err)
		}

		return nil, fmt.Errorf("unable to decode jpeg: %w", err)
	}

	var side1, side2 int

	if img.Bounds().Dx() > img.Bounds().Dy() {
		side1 = size
	} else {
		side2 = size
	}

	resized := imaging.Resize(img, side1, side2, imaging.Lanczos)

	// encode resized image to png
	buf := new(bytes.Buffer)

	err = png.Encode(buf, resized)
	if err != nil {
		return nil, fmt.Errorf("unable to encode resised png: %w", err)
	}

	return buf.Bytes(), nil
}
