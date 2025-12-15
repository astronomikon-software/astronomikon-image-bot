package postimgs

import (
	"errors"
	"regexp"
)

var ErrLinkNotFound = errors.New("image link not found")

func ExtractImageUrl(html string) (string, error) {
	pattern := regexp.MustCompile("https://cdn\\.donmai\\.us/original/[\\w/]*\\.[\\wd]*")
	match := pattern.FindStringSubmatch(html)
	if len(match) == 0 {
		return "", ErrLinkNotFound
	} else {
		return findLongestItem(match), nil
	}
}

func findLongestItem(items []string) string {
	if len(items) == 0 {
		return ""
	}
	bestIndex := 0
	maxLen := len(items[0])
	for i, item := range items {
		curLen := len(item)
		if curLen > maxLen {
			bestIndex = i
			maxLen = curLen
		}
	}
	return items[bestIndex]
}
