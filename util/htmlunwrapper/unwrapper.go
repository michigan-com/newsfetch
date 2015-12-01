package htmlunwrapper

import (
	"golang.org/x/net/html"

	"github.com/andybalholm/cascadia"
	"github.com/michigan-com/newsfetch/util/htmldomutil"
	"github.com/michigan-com/newsfetch/util/orderedlist"
)

type Tag string

type Unwrapper struct {
	// parallel arrays
	selectors []cascadia.Selector
	tags      []string
}

func NewFromTable(m map[string]string) *Unwrapper {
	u := new(Unwrapper)
	u.AddFromTable(m)
	return u
}

// AddFromTable accepts a map of selectors to tags.
func (u *Unwrapper) AddFromTable(m map[string]string) {
	for selector, tag := range m {
		u.Add(selector, tag)
	}
}

func (u *Unwrapper) Add(selector string, tag string) {
	sel := cascadia.MustCompile(selector)

	u.selectors = append(u.selectors, sel)
	u.tags = append(u.tags, tag)
}

func (u *Unwrapper) Analyze(node *html.Node) []string {
	_, tags := u.Unwrap(node)
	return tags
}

func (u *Unwrapper) Unwrap(node *html.Node) (*html.Node, []string) {
	var tags []string

	for {
		tags = u.analyze(node, tags)

		child := htmldomutil.FindSoleChildElement(node)
		if child == nil {
			break
		} else {
			node = child
		}
	}

	return node, tags
}

func (u *Unwrapper) analyze(node *html.Node, tags []string) []string {
	for i, selector := range u.selectors {
		if selector.Match(node) {
			tag := u.tags[i]
			tags = orderedlist.InsertString(tags, tag, true)
		}
	}
	return tags
}
