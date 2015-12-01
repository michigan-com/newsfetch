package htmldomutil

import (
	"strings"

	"golang.org/x/net/html"
)

func ParseString(htmlstr string) (*html.Node, error) {
	return html.Parse(strings.NewReader(htmlstr))
}

func MustParseString(htmlstr string) *html.Node {
	doc, err := ParseString(htmlstr)
	if err != nil {
		panic(err)
	}
	return doc
}
