package htmldomutil

import (
	"bytes"

	"golang.org/x/net/html"
)

// based on https://github.com/PuerkitoBio/goquery/blob/master/property.go

func GetNodeText(node *html.Node) string {
	buf := new(bytes.Buffer)
	collectNodeText(node, buf)
	return buf.String()
}

func collectNodeText(node *html.Node, buf *bytes.Buffer) {
	if node.Type == html.TextNode {
		// Keep newlines and spaces, like jQuery
		buf.WriteString(node.Data)

	} else if node.FirstChild != nil {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			collectNodeText(c, buf)
		}
	}
}
