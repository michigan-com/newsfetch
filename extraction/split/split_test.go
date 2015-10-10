package split

import (
	"strings"
	"testing"
)

func TestSplitSimple(t *testing.T) {
	os(t, "Hello, world!", "Hello / world")
}

func TestSplitTwitter(t *testing.T) {
	os(t, "Contact @foo.", "Contact / @foo")
	os(t, "Tweet with #CoolHashtag!", "Tweet / with / #CoolHashtag")
}

func TestSplitEmail(t *testing.T) {
	os(t, "Contact someone@example.com.", "Contact / someone@example.com")
}

func TestSplitURL(t *testing.T) {
	os(t, "Contact http://example.com/.", "Contact / http://example.com/")
}

func TestSplitPhone(t *testing.T) {
	os(t, "Contact 555-55-55.", "Contact / 555-55-55")
}

func os(t *testing.T, input string, expected string) {
	actual := strings.Join(SplitWords(input), " / ")
	if actual != expected {
		t.Errorf("Split(%#v) != %#v", actual, expected)
	}
}
