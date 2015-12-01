package htmldomutil

import (
	"golang.org/x/net/html"

	"github.com/michigan-com/newsfetch/util/stringutil"
)

func FixupTextNodesBeforeSubSup(parent *html.Node) {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if (child.Type == html.ElementNode) && (child.Data == "sup") {
			if prev := child.PrevSibling; prev != nil && prev.Type == html.TextNode {
				if stringutil.EndsWithDigit(prev.Data) {
					prev.Data = prev.Data + " "
				}
			}

		}
		if child.Type == html.ElementNode {
			FixupTextNodesBeforeSubSup(child)
		}
	}
}
