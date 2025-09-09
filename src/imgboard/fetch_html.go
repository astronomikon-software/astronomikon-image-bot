package imgboard

import (
	"errors"
	"io"
	"net/http"
)

var ErrFetchFailed = errors.New("post image: failed to fetch page")

func FetchHtml(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return "", ErrFetchFailed
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()
	return string(body), nil
}
