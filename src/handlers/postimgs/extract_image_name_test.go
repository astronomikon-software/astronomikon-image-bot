package postimgs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractImageName(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedName string
		expectedErr  error
	}{
		{
			name:         "empty url",
			url:          "",
			expectedName: "",
			expectedErr:  ErrInvalidImageUrl,
		},
		{
			name:         "invalid url",
			url:          "url",
			expectedName: "",
			expectedErr:  ErrInvalidImageUrl,
		},
		{
			name:         "valid url",
			url:          "https://cdn.donmai.us/original/15/d9/image.png",
			expectedName: "image.png",
			expectedErr:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, err := ExtractImageName(test.url)
			assert.Equal(t, test.expectedName, name)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
