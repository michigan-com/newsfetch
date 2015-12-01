package htmltotext

import (
	"golang.org/x/net/html"

	"github.com/michigan-com/newsfetch/util/htmldomutil"
	"github.com/michigan-com/newsfetch/util/htmlunwrapper"
	"github.com/michigan-com/newsfetch/util/stringutil"

	t "github.com/michigan-com/newsfetch/model/bodytypes"
)

var unwrapper = htmlunwrapper.NewFromTable(map[string]string{
	".-newsgate-paragraph-cci-howto-head-":       t.NewsgateHead,
	".-newsgate-paragraph-cci-howto-components-": t.NewsgateComponent,
	".-newsgate-element-cci-howto--end":          t.NewsgateEnd,
	"li":        t.ListItem,
	"b, strong": t.Bold,
})

func ConvertParagraphNode(node *html.Node) t.Paragraph {
	htmldomutil.FixupTextNodesBeforeSubSup(node)

	text := stringutil.NormalizeSpace(htmldomutil.GetNodeText(node))

	tags := unwrapper.Analyze(node)

	return t.Paragraph{text, tags}
}
