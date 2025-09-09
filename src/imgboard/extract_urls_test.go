package imgboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractUrls(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		expectedUrls []string
		expectedErr  error
	}{
		{
			name:         "empty message",
			message:      "",
			expectedUrls: []string{},
			expectedErr:  nil,
		},
		{
			name:         "single line message",
			message:      "https://betabooru.donmai.us/posts/123456789",
			expectedUrls: []string{"https://betabooru.donmai.us/posts/123456789"},
			expectedErr:  nil,
		},
		{
			name: "multiline line message",
			message: "https://betabooru.donmai.us/posts/123456789\n" +
				"   https://betabooru.donmai.us/posts/234567891\n" +
				"https://betabooru.dOnmai.US/posts/345678912?q=xyz_123  \n",
			expectedUrls: []string{
				"https://betabooru.donmai.us/posts/123456789",
				"https://betabooru.donmai.us/posts/234567891",
				"https://betabooru.donmai.us/posts/345678912",
			},
			expectedErr: nil,
		},
		{
			name: "invalid url message",
			message: "https://betabooru.donmai.us/posts/123456789\n" +
				"   https://betabooru.donmai.us/posts/",
			expectedUrls: nil,
			expectedErr:  ErrInvalidPageUrl,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls, err := ExtractUrls(test.message)
			assert.ElementsMatch(t, test.expectedUrls, urls)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
