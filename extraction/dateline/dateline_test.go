package dateline

import (
	"testing"
)

func TestRemoveDatelineUppercase(t *testing.T) {
	o(t, "Hello, world!", "Hello, world!")
	o(t, "WASHINGTON — Hello, world!", "Hello, world!")
	o(t, "WASHINGTON — 'Hello, world!'", "'Hello, world!'")
	o(t, "WASHINGTON — “Hello, world!”", "“Hello, world!”")
	o(t, "WASHINGTON -- Hello, world!", "Hello, world!")
	o(t, "ORLANDO, FLA. -- Hello, world!", "Hello, world!")
}

func TestRemoveDatelineAlternative(t *testing.T) {
	o(t, "Orlando, Fla. -- Hello, world!", "Hello, world!")
	o(t, "Some unrelated text which is not a dateline -- Hello, world!", "Some unrelated text which is not a dateline -- Hello, world!")
}

func o(t *testing.T, input string, expected string) {
	actual := RemoveDateline(input)
	if actual != expected {
		t.Errorf("For %#v expected %#v, got %#v", input, expected, actual)
	}
}
