package goqueryutil

import (
	gq "github.com/PuerkitoBio/goquery"
	"github.com/michigan-com/newsfetch/util/htmldomutil"
)

func FixupTextNodesBeforeSubSup(s *gq.Selection) {
	for _, node := range s.Nodes {
		htmldomutil.FixupTextNodesBeforeSubSup(node)
	}
}
