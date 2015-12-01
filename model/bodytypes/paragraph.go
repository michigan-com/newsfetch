package bodytypes

import (
	"bytes"

	"github.com/michigan-com/newsfetch/util/orderedlist"
)

const (
	NewsgateHead      = "newsgate-head"
	NewsgateComponent = "newsgate-component"
	NewsgateEnd       = "newsgate-end"
	ListItem          = "list"
	Bold              = "bold"
)

type Paragraph struct {
	Text string
	Tags []string
}

func (p Paragraph) String() string {
	var buf bytes.Buffer
	for _, tag := range p.Tags {
		buf.WriteString("#")
		buf.WriteString(tag)
		buf.WriteString(" ")
	}
	buf.WriteString(p.Text)
	return buf.String()
}

func (p Paragraph) HasTag(tag string) bool {
	return orderedlist.ContainsString(p.Tags, tag)
}
