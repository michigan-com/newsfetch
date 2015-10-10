package classify

import (
	"testing"
)

func TestClassifyValid(t *testing.T) {
	oc(t, "Hello, world!", true)
	oc(t, `"We've got a lot of work to do," McCabe said.`, true)
	oc(t, `Tigers general manager Al Avila today said that no decision has been made.`, true)
}

func TestClassifyTrailingLines(t *testing.T) {
	oc(t, "Follow him on Twitter: @anthonyfenech.", false)
	oc(t, "Check out our latest Tigers podcast at freep.com/tigerspodcast or on iTunes.", false)
	oc(t, "Contact L.L. Brasier: 248-858-2262 or lbrasier@freepress.com", false)
	oc(t, "And download our free Tigers Xtra app on Apple and Android!", false)
	oc(t, "Contact Daniel Bethencourt: dbethencourt@freepress.com or 313-223-4531. Follow on Twitter at @_dbethencourt.", false)
	oc(t, "Anyone with details on what happened can call the Auburn Hills Police Department at 248-370-9444.", false)
	oc(t, "Follow him on Twitter: @anthonyfenech. Check out our latest Tigers podcast at freep.com/tigerspodcast or on iTunes. And download our free Tigers Xtra app on Apple and Android!", false)
}

func TestClassifyPoliceLines(t *testing.T) {
	oc(t, "Anyone with details on what happened can call the Auburn Hills Police Department at 248-370-9444.", false)
}

func oc(t *testing.T, input string, expected bool) {
	actual, rationale := IsWorthyParagraph(input)
	if actual != expected {
		if actual {
			t.Errorf("Errorneously classified as GOOD: %#v", input)
		} else {
			t.Errorf("Errorneously classified as BAD: %#v", input)
		}
		t.Logf("Decision based on: %s", rationale)
	}
}
