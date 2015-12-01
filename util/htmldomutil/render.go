package htmldomutil

import (
	"bytes"

	"golang.org/x/net/html"
)

func RenderToString(node *html.Node) string {
	buf := new(bytes.Buffer)
	err := html.Render(buf, node)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
