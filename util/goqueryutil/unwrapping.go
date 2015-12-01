package goqueryutil

import (
	gq "github.com/PuerkitoBio/goquery"
	"github.com/michigan-com/newsfetch/util/htmldomutil"
	// "golang.org/x/net/html"
)

func HasSingleChildMatching(s *gq.Selection, selector string) bool {
	if htmldomutil.FindSoleChildElement(s.Nodes[0]) == nil {
		return false
	}

	return s.Children().Is(selector)
}
