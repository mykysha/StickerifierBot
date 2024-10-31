package photoconverter

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
)

const tgPhotoSize = 512

var errUnsupportedFormat = errors.New("unsupported photo format")

func ConvertToPNG(photo []byte) ([]byte, error) {
	resized, err := resizePNG(photo, tgPhotoSize)
	if err != nil {
		return nil, fmt.Errorf("unable to resize png: %w", err)
	}

	return resized, nil
}

func resizePNG(photo []byte, size int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(photo))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, fmt.Errorf("%w: %w", errUnsupportedFormat, err)
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

	buf := new(bytes.Buffer)

	if err := png.Encode(buf, resized); err != nil {
		return nil, fmt.Errorf("unable to encode resised png: %w", err)
	}

	return buf.Bytes(), nil
}
