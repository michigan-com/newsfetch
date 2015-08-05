package lib

import (
	"net/url"
	"regexp"
	"strconv"
)

func GetArticleId(url string) int {
	// Given an article url, get the ID from it
	r := regexp.MustCompile("/([0-9]+)/{0,1}$")
	match := r.FindStringSubmatch(url)

	if len(match) > 1 {
		i, err := strconv.Atoi(match[1])
		if err != nil {
			return -1
		}
		return i
	} else {
		return -1
	}
}

func GetHost(inputUrl string) (string, error) {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}

	replace := regexp.MustCompile("[.].+$")
	return replace.ReplaceAllString(u.Host, ""), nil
}
