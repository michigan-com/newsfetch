package fuzzy_classifier

import (
	"strings"
	"testing"

	"github.com/michigan-com/newsfetch/extraction/diff"
)

func TestClassifierBuiltInTags(t *testing.T) {
	classifier := NewFuzzyClassifier()

	assertClassificationResult(t, classifier, `
        Follow him on Twitter: @example.
    `, `
        @cap Follow
        @s him
        @s on
        @cap Twitter:
        @twitter @s @example.
    `)

	assertClassificationResult(t, classifier, `
        Add 1 1/2 cups of baking flour.
    `, `
        @cap Add
        @number 1 1/2
        @integer 1
        @fraction 1/2
        cups
        @s of
        baking
        flour.
    `)
}

func TestClassifierTrivialTag(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@foo
        bar boz
        fubar
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        xxx bar bar boz boz yyy fubar zzz
    `, `
        xxx
        bar
        @foo bar boz
        bar
        boz
        boz
        yyy
        @foo fubar
        zzz
    `)
}

func TestClassifierSubtag(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@follow
        follow @someone on twitter

        :@someone
        him
        her
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        Follow him on Twitter.
    `, `
        @follow Follow him on Twitter.
        @cap Follow
        @someone @s him
        @s on
        @cap Twitter.
    `)
}

func TestClassifierOptional(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@follow
        follow ?@someone on twitter

        :@someone
        him
        her
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        Follow on Twitter.
    `, `
        @follow Follow on Twitter.
        @cap Follow
        @s on
        @cap Twitter.
    `)
}

func TestClassifierOptional2(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@ingredient
        some ?@color @noun

        :@color
        black

        :@noun
        pepper
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        some black pepper
    `, `
        @ingredient some black pepper
        @s some
        @color black
        @noun pepper
    `)
	assertClassificationResult(t, classifier, `
        some pepper
    `, `
        @ingredient some pepper
        @s some
        @noun pepper
    `)
}

func TestClassifierOptional3(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@ingredient
        ?@color @noun

        :@color
        black

        :@noun
        pepper
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        some black pepper
    `, `
        @s some
        @ingredient black pepper
        @color black
        @noun pepper
    `)
	assertClassificationResult(t, classifier, `
        some pepper
    `, `
        @s some
        @ingredient @noun pepper
    `)
}

func TestClassifierLeadingOptional(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@ingredient
        ?@color @noun

        :@color
        black

        :@noun
        pepper
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        black pepper
    `, `
        @ingredient black pepper
        @color black
        @noun pepper
    `)
	assertClassificationResult(t, classifier, `
        pepper
    `, `
        @ingredient @noun pepper
    `)
}

func TestClassifierSkipBuiltInWords(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@follow
        .skip @s
        follow twitter
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        Follow him on Twitter.
    `, `
        @follow Follow him on Twitter.
        @cap Follow
        @s him
        @s on
        @cap Twitter.
    `)
}

func TestClassifierSkipCustomSkippable(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@follow
        .skip @s @publication
        follow twitter

        :@publication
        detroit free press
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        Follow Detroit Free Press on Twitter.
    `, `
        @follow Follow Detroit Free Press on Twitter.
        @cap Follow
        @publication Detroit Free Press
        @cap Detroit
        @cap Free
        @cap Press
        @s on
        @cap Twitter.
    `)
}

func TestClassifierSkippablePrefix(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@follow
        .skip +b @s @pleasantries @pronoun
        follow twitter

        :@pleasantries
        please

        :@pronoun
        us
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        Please, please, please follow us on Twitter!
    `, `
        @follow Please, please, please follow us on Twitter!
        @pleasantries @cap Please,
        @pleasantries please,
        @pleasantries please
        follow
        @pronoun us
        @s on
        @cap Twitter!
    `)
}

func TestClassifierSkippableSuffix(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@ingredient
        .skip +b +a @adjective
        @noun

        :@adjective
        freshly ground
        organic
        black

        :@noun
        pepper
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        Add some black pepper, freshly ground, organic.
    `, `
        @cap Add
        @s some
        @ingredient black pepper, freshly ground, organic.
        @adjective black
        @noun pepper,
        @adjective freshly ground,
        freshly
        ground,
        @adjective organic.
    `)
}

func TestClassifierSkippableSuffix2(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@ingredient
        .skip +b +a @alternative @purpose_clause
        @noun

        :@alternative
        or another herb

        :@purpose_clause
        for garnish

        :@noun
        a few sprigs of flat-leaf parsley
        herb
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        A few sprigs of flat-leaf parsley (or another herb) for garnish
    `, `
        @ingredient A few sprigs of flat-leaf parsley (or another herb) for garnish
        @noun A few sprigs of flat-leaf parsley
        @s A
        @s few
        sprigs
        @s of
        @s flat-leaf
        parsley
        @alternative (or another herb)
        @s (or
        another
        @noun herb)
        @purpose_clause for garnish
        @s for
        garnish
    `)
}

func TestClassifierSkipJustTheRightAmount(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@follow
        .skip @pleasantries @publication @pronoun
        please follow @command on twitter

        :@pleasantries
        would you please
        please
        pretty please

        :@pronoun
        us

        :@publication
        us press
        detroit free press

        :@command
        press it now
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        Would you please follow US (press it now!) (pretty please!) on Twitter?
    `, `
        @pleasantries Would you please
        @cap Would
        @s you
        @follow please follow US (press it now!) (pretty please!) on Twitter?
        please
        follow
        @publication US (press
        @pronoun US
        @command (press it now!)
        (press
        @s it
        @s now!)
        @pleasantries (pretty please!)
        (pretty
        please!)
        @s on
        @cap Twitter?
    `)
}

func TestClassifierRepetition(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@colors
        .skip @or
        +@color

        :@color
        red
        green
        blue

        :@or
        or
        and
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        red, green and blue and
    `, `
        @colors red, green and blue
        @color red,
        @color green
        @or @s and
        @color blue
        @or @s and
    `)
}

func TestClassifierLeadingRepetition(t *testing.T) {
	classifier := NewFuzzyClassifier()
	err := classifier.Add(`
        :@ingredient
        @parts of @noun
        @noun

        :@parts
        .skip @or
        +@part

        :@part
        curls
        julienne

        :@or
        or
        and

        :@noun
        lemon zest
    `)
	if err != nil {
		t.Fatalf("Classifier returned a parse error: %v", err.Error())
	}

	assertClassificationResult(t, classifier, `
        curls or julienne of lemon zest
    `, `   
        @ingredient curls or julienne of lemon zest
        @parts curls or julienne
        @part curls
        @or @s or
        @part julienne
        @s of
        @noun lemon zest
        lemon
        zest
    `)
}

func assertClassificationResult(t *testing.T, c *Classifier, input string, expected string) {
	result := c.Process(input)

	actual := diff.TrimLinesInString(result.Description())
	expected = diff.TrimLinesInString(expected)

	if actual != expected {
		t.Errorf("Classification result mismatch for %#v.", strings.TrimSpace(input))
		t.Logf("Diff:\n%v", diff.LineDiff(expected, actual))
		// t.Logf("actual %#v != expected %#v", actual, expected)
		t.Log("------------------------------")
		t.Logf("Actual:\n%v", actual)
		t.Log("------------------------------")
		t.Logf("Expected:\n%v", expected)
	}
}
