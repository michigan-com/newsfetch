package fuzzy_classifier

import (
	"strings"
	"testing"

	"github.com/michigan-com/newsfetch/extraction/diff"
)

func TestClassifierParseErrors(t *testing.T) {
	assertParseError(t, `
        foo
    `, "expected a category header")
	assertParseError(t, `
        :foo
    `, "expected @ after :")
	assertParseError(t, `
        :@ foo
    `, "expected a tag name after :@")
	assertParseError(t, `
        :@foo
        .bar
    `, "unknown instruction \".bar\"")
	assertParseError(t, `
        :@foo
        .skip bar
    `, "expected @")
}

func TestClassifierParseSimpleSpec(t *testing.T) {
	assertParseResult(t, `
        :@follow
        follow @someone on twitter

        :@someone
        him
        her
    `, `
        CAT @follow
        SCHEME follow <@someone> on twitter

        CAT @someone
        SCHEME him
        SCHEME her
    `)
}

func TestClassifierParseAttributes(t *testing.T) {
	assertParseResult(t, `
        :@follow
        .skip @someone @something
        follow on twitter
    `, `
        CAT @follow
        SKIP @someone @something
        SCHEME follow on twitter
    `)
}

func assertParseError(t *testing.T, definition string, message string) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(definition)
	if err == nil {
		t.Fatalf("Classifier didn't produce a parse error")
	} else if !strings.Contains(err.Error(), message) {
		t.Fatalf("Classifier parse error %+v does not contain expected string %+v", err.Error(), message)
	}
}

func assertParseResult(t *testing.T, definition string, expected string) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(definition)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	actual := diff.TrimLinesInString(classifier.Description())
	expected = diff.TrimLinesInString(expected)

	if actual != expected {
		t.Errorf("Parse result mismatch.")
		t.Logf("Diff:\n%v", diff.LineDiff(expected, actual))
		// t.Logf("actual %#v != expected %#v", actual, expected)
		t.FailNow()
	}
}
