package postimgs

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidPageUrl = errors.New("invalid page url")
)

func ExtractUrls(message string) ([]string, error) {
	var res []string
	for _, line := range strings.Split(message, "\n") {
		if !isBlank(line) {
			url, err := extractValidUrl(line)
			if err != nil {
				return nil, err
			}
			res = append(res, url)
		}
	}
	return res, nil
}

func isBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func extractValidUrl(s string) (string, error) {
	for _, pattern := range []*regexp.Regexp{
		regexp.MustCompile("^https://betabooru\\.donmai\\.us/posts/\\d+"),
		regexp.MustCompile("^https://danbooru\\.donmai\\.us/posts/\\d+"),
	} {
		res := pattern.FindString(strings.ToLower(strings.TrimSpace(s)))
		if res != "" {
			return res, nil
		}
	}
	return "", fmt.Errorf("%v: %w", s, ErrInvalidPageUrl)
}
