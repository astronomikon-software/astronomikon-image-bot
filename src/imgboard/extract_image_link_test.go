package imgboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractImageLink(t *testing.T) {
	tests := []struct {
		name         string
		html         string
		expectedLink string
		expectedErr  error
	}{
		{
			name:         "empty html",
			html:         "",
			expectedLink: "",
			expectedErr:  ErrLinkNotFound,
		},
		{
			name:         "missing link html",
			html:         "<a class=\"image-view-original-link\">view original</a>",
			expectedLink: "",
			expectedErr:  ErrLinkNotFound,
		},
		{
			name:         "single tag",
			html:         "<a class=\"image-view-original-link\" href=\"https://cdn.donmai.us/original/15/d9/image.png\">view original</a>",
			expectedLink: "https://cdn.donmai.us/original/15/d9/image.png",
			expectedErr:  nil,
		},
		{
			name: "multiple tags",
			html: "<a class=\"image-view-original-link\" href=\"https://cdn.donmai.us/original/15/d9/image0.png\">view original</a>\n" +
				"<a class=\"image-view-original-link\" href=\"https://cdn.donmai.us/original/15/d9/image1.png\">view original</a>",
			expectedLink: "https://cdn.donmai.us/original/15/d9/image0.png",
			expectedErr:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			link, err := ExtractImageUrl(test.html)
			assert.Equal(t, test.expectedLink, link)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
