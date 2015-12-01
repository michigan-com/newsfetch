package fuzzy

import (
	"testing"
)

func TestClassifyNumber(t *testing.T) {
	onum(t, "", NumClassNone)
	onum(t, "abc", NumClassNone)

	onum(t, "1", NumClassInteger)
	onum(t, "123", NumClassInteger)

	onum(t, "123.", NumClassFloat)
	onum(t, "123.4", NumClassFloat)
	onum(t, "123.456", NumClassFloat)
	onum(t, "1.456", NumClassFloat)
	onum(t, ".456", NumClassFloat)
	onum(t, ".456...", NumClassFloat)

	onum(t, "1/2", NumClassFraction)
	onum(t, "12/34", NumClassFraction)
	onum(t, "/34", NumClassNone)
	onum(t, "1/", NumClassNone)
	onum(t, "12/34...", NumClassFraction)

	onum(t, "$12", NumClassCurrency)
	onum(t, "$12.45", NumClassCurrency)
	onum(t, "$12k", NumClassCurrency)
	onum(t, "$12m", NumClassCurrency)
	onum(t, "$12MM", NumClassCurrency)
	onum(t, "$12b", NumClassCurrency)
	onum(t, "$12bb", NumClassCurrency)
	onum(t, "$12bbb", NumClassNone)
	onum(t, "$12bb...", NumClassCurrency)
}

func onum(t *testing.T, input string, expected NumClass) {
	actual := ClassifyNumberString(input)

	if actual != expected {
		t.Errorf("%#v is %d, expected %d.", input, actual, expected)
	}
}
