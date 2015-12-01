package htmltotext

import (
	"testing"

	"github.com/andybalholm/cascadia"
	"github.com/michigan-com/newsfetch/util/htmldomutil"
)

func TestConvertParagraph(t *testing.T) {
	op(t, `<p>hello</p>`, "hello")

	op(t, `<p><b>hello</b></p>`, "#bold hello")
	op(t, `<p><strong>hello</strong></p>`, "#bold hello")

	op(t, `<li>hello</li>`, "#list hello")
	op(t, `<li><b>hello</b></li>`, "#bold #list hello")

	op(t, `<div class="-newsgate-paragraph-cci-howto-head-">hello</div>`, "#newsgate-head hello")
	op(t, `<div class="-newsgate-paragraph-cci-howto-components-">hello</div>`, "#newsgate-component hello")
	op(t, `<div class="-newsgate-element-cci-howto--end">hello</div>`, "#newsgate-end hello")
}

func op(t *testing.T, input string, expected string) {
	root := htmldomutil.MustParseString(input)
	inputEl := cascadia.MustCompile("body").MatchFirst(root).FirstChild

	para := ConvertParagraphNode(inputEl)

	actual := para.String()

	if actual != expected {
		inputHTML := htmldomutil.RenderToString(inputEl)
		t.Errorf("Convert(%v) = %v; expected: %v", inputHTML, actual, expected)
	}
}
