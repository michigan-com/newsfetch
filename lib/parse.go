package lib

import (
	"errors"
	"fmt"
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

/*
	Get the url host from the url string (inputUrl)

	Ex:
		result, err := GetHost("http://google.com")
		// result == "google"

	Using the url.Parse method, so urls must start with "http://"

*/
func GetHost(inputUrl string) (string, error) {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}

	hostRegex := regexp.MustCompile("([a-zA-Z0-9][a-zA-Z0-9-_]{0,61}[a-zA-Z0-9]{0,1}).[a-zA-Z]{2,}$")
	match := hostRegex.FindStringSubmatch(u.Host)
	if match == nil {
		return "", errors.New(fmt.Sprintf("Could not get host from %s", u.Host))
	}

	return match[1], nil
}
