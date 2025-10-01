package imgproc

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math"

	"github.com/disintegration/imaging"
)

func ResizeImage(imgBytes []byte, maxSum int, maxSize int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	scale := math.Min(
		ScaleFactorByDimensionsSum(width, height, maxSum),
		ScaleFactorBySize(len(imgBytes), maxSize),
	)

	if scale != 1 {
		return ScaleImage(img, scale)
	} else {
		return imgBytes, nil
	}
}

func ScaleImage(img image.Image, scaleFactor float64) ([]byte, error) {
	newImage := imaging.Resize(
		img,
		int(float64(img.Bounds().Dx())*scaleFactor),
		int(float64(img.Bounds().Dy())*scaleFactor),
		imaging.MitchellNetravali,
	)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, newImage, nil)
	if err != nil {
		return nil, fmt.Errorf("error encoding image: %w", err)
	}
	return buf.Bytes(), nil
}

func ScaleFactorByDimensionsSum(width int, height int, maxSum int) float64 {
	currentSum := width + height
	if currentSum > maxSum {
		return float64(maxSum) / float64(currentSum)
	} else {
		return 1
	}
}

func ScaleFactorBySize(currentSize int, maxSize int) float64 {
	if currentSize > maxSize {
		return math.Sqrt(float64(maxSize) / float64(currentSize))
	} else {
		return 1
	}
}
