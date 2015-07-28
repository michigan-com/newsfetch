package lib

import (
	"net/url"
	"regexp"
)

func GetHost(inputUrl string) string {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return ""
	}

	replace := regexp.MustCompile("[.].+$")
	return replace.ReplaceAllString(u.Host, "")
}
