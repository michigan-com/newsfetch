package htmldomutil

import (
	"golang.org/x/net/html"
)

func FindSoleChildElement(parent *html.Node) *html.Node {
	var result *html.Node

	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		switch child.Type {
		case html.CommentNode:
		case html.TextNode:
			if child.Data != "" {
				return nil
			}
		case html.ElementNode:
			if result != nil {
				return nil // found second child
			} else {
				result = child
			}
		default:
			return nil
		}
	}

	return result
}
