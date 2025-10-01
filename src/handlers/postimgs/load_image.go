package postimgs

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Image struct {
	Url      string
	Filename string
	Data     []byte
}

func LoadImage(ctx context.Context, pageUrl string) (*Image, error) {
	logFetching(pageUrl)
	html, err := FetchHtml(ctx, pageUrl)
	if err != nil {
		return nil, err
	}

	imageUrl, err := ExtractImageUrl(html)
	if err != nil {
		return nil, err
	}

	imageName, err := ExtractImageName(imageUrl)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Image{
		Url:      imageUrl,
		Filename: imageName,
		Data:     data,
	}, nil
}

func logFetching(url string) {
	fmt.Printf("%s: fetching %q\n", logPrefix, url)
}
