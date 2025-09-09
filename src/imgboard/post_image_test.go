package imgboard

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostImage(t *testing.T) {
	testCases := []struct {
		name    string
		pageUrl string

		fetcherCalls int
		fetcherRes   string
		fetcherErr   error

		senderCalls int
		senderImage Image
		senderErr   error

		expectedErr error
	}{
		{
			name:    "Success",
			pageUrl: "https://page/url",

			fetcherCalls: 1,
			fetcherRes:   "<a class=\"image-view-original-link\" href=\"https://cdn.donmai.us/original/15/d9/image.png\">view original</a>",
			fetcherErr:   nil,

			senderCalls: 1,
			senderErr:   nil,
			senderImage: Image{Url: "https://cdn.donmai.us/original/15/d9/image.png", Filename: "image.png"},

			expectedErr: nil,
		},
		{
			name:    "Fetch failed",
			pageUrl: "https://page/url",

			fetcherCalls: 1,
			fetcherRes:   "",
			fetcherErr:   ErrFetchFailed,

			senderCalls: 0,
			senderErr:   nil,
			senderImage: Image{},

			expectedErr: ErrFetchFailed,
		},
		{
			name:    "Extract Url failed",
			pageUrl: "https://page/url",

			fetcherCalls: 1,
			fetcherRes:   "<a class=\"image-view-original-link\">view original</a>",
			fetcherErr:   nil,

			senderCalls: 0,
			senderErr:   nil,
			senderImage: Image{},

			expectedErr: ErrLinkNotFound,
		},
		{
			name:    "Extract image name failed",
			pageUrl: "https://page/url",

			fetcherCalls: 1,
			fetcherRes:   "<a class=\"image-view-original-link\" href=\"invalid-Url\">view original</a>",
			fetcherErr:   nil,

			senderCalls: 0,
			senderErr:   nil,
			senderImage: Image{},

			expectedErr: ErrInvalidImageUrl,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.TODO()

			fetcherCalls := 0
			fetcherMock := func(pageUrl string) (string, error) {
				fetcherCalls++
				assert.Equal(t, testCase.pageUrl, pageUrl)
				return testCase.fetcherRes, testCase.fetcherErr
			}

			senderCalls := 0
			senderMock := func(ctx context.Context, image Image) error {
				senderCalls++
				assert.Equal(t, testCase.senderImage.Url, image.Url)
				assert.Equal(t, testCase.senderImage.Filename, image.Filename)
				return testCase.fetcherErr
			}

			err := PostImage(ctx, testCase.pageUrl, fetcherMock, senderMock)
			assert.Equal(t, 1, fetcherCalls)
			assert.Equal(t, testCase.senderCalls, senderCalls)
			assert.ErrorIs(t, err, testCase.expectedErr)
		})
	}
}
