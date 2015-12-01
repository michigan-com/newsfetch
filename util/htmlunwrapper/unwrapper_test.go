package htmlunwrapper

import (
	"strings"
	"testing"

	"github.com/andybalholm/cascadia"
	"github.com/michigan-com/newsfetch/util/htmldomutil"
)

func TestUnwrapping(t *testing.T) {
	u := NewFromTable(map[string]string{
		".container": "#container",
		"b, strong":  "#bold",
	})

	o(t, u, "<p>hello</p>", "")
	o(t, u, "<p><b>hello</b></p>", "#bold")
	o(t, u, "<p><b><strong>hello</strong></b></p>", "#bold")
	o(t, u, `<p><span class="container"><strong>hello</strong></span></p>`, "#bold #container")
	o(t, u, `<p><b><span class="container">hello</span></b></p>`, "#bold #container")
	o(t, u, `<div class="container"><p><strong>hello</strong></p></div>`, "#bold #container")
}

func o(t *testing.T, u *Unwrapper, input string, expected string) {
	root := htmldomutil.MustParseString(input)
	inputEl := cascadia.MustCompile("body").MatchFirst(root).FirstChild

	tags := u.Analyze(inputEl)

	inputHTML := htmldomutil.RenderToString(inputEl)

	actual := strings.TrimSpace(strings.Join(tags, " "))

	if actual != expected {
		t.Errorf("Unwrap(%v) = %v; expected: %v", inputHTML, actual, expected)
	}
}
