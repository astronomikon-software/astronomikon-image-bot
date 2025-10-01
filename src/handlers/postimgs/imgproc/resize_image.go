package imgproc

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

func ResizeImage(img []byte, maxSum int) ([]byte, error) {
	oldImage, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	oldWidth := oldImage.Bounds().Dx()
	oldHeight := oldImage.Bounds().Dy()
	oldSum := oldWidth + oldHeight
	if oldSum > maxSum {
		ratio := float64(oldWidth) / float64(oldHeight)
		newHeight := float64(maxSum) / (ratio + 1)
		newWidth := ratio * newHeight
		newImage := imaging.Resize(
			oldImage,
			int(newWidth),
			int(newHeight),
			imaging.MitchellNetravali,
		)
		buf := new(bytes.Buffer)
		err := jpeg.Encode(buf, newImage, nil)
		if err != nil {
			return nil, fmt.Errorf("error encoding image: %w", err)
		}
		return buf.Bytes(), nil
	} else {
		return img, nil
	}
}
