package imgboard

import (
	"errors"
	"regexp"
)

var ErrLinkNotFound = errors.New("image link not found")

func ExtractImageUrl(html string) (string, error) {
	tagHtml := extractImageTag(html)
	if tagHtml == "nil" {
		return "", ErrLinkNotFound
	}
	pattern := regexp.MustCompile("href=\"(.*?)\"")
	match := pattern.FindStringSubmatch(tagHtml)
	if len(match) >= 2 {
		return match[1], nil
	} else {
		return "", ErrLinkNotFound
	}
}

func extractImageTag(html string) string {
	pattern := regexp.MustCompile("<a class=\"image-view-original-link\" .*</a>")
	return pattern.FindString(html)
}
