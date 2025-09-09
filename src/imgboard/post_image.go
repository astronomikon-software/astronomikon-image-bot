package imgboard

import (
	"context"
)

type Image struct {
	Url      string
	Filename string
}

type SenderFunc func(ctx context.Context, image Image) error

type FetcherFunc func(pageUrl string) (string, error)

func PostImage(ctx context.Context, pageUrl string, fetcher FetcherFunc, sender SenderFunc) error {
	html, err := fetcher(pageUrl)
	if err != nil {
		return err
	}

	imageUrl, err := ExtractImageUrl(html)
	if err != nil {
		return err
	}

	imageName, err := ExtractImageName(imageUrl)
	if err != nil {
		return err
	}

	return sender(ctx, Image{
		Url:      imageUrl,
		Filename: imageName,
	})
}
