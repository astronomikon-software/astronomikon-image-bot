package imgboard

import (
	"errors"
	"strings"
)

var ErrInvalidImageUrl = errors.New("invalid image url")

func ExtractImageName(url string) (string, error) {
	tokens := strings.Split(url, "/")
	if len(tokens) >= 2 {
		return tokens[len(tokens)-1], nil
	} else {
		return "", ErrInvalidImageUrl
	}
}
