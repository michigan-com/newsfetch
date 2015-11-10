package fuzzy_classifier

import (
	"unicode"
)

type NumClass int

const (
	NumClassNone NumClass = iota
	NumClassInteger
	NumClassIntegerWithMagniture
	NumClassCurrency
	NumClassFloat
	NumClassFraction
)

func ClassifyNumberString(s string) NumClass {
	return ClassifyNumber([]rune(s))
}

func ClassifyNumber(runes []rune) NumClass {
	isCurrency := false
	var magnitudeSpec rune
	state := 0
	inTrailing := false
	for _, r := range runes {
		if inTrailing {
			if unicode.IsPunct(r) {
				continue
			} else {
				return NumClassNone
			}
		}
		switch state {
		case 0: // start
			if r == '$' || r == '€' {
				isCurrency = true
				state = 4
			} else if unicode.IsDigit(r) {
				state = 1
			} else if r == '.' {
				state = 5
			} else if unicode.IsPunct(r) {
				return NumClassNone
			}
		case 1: // within a run of leading digits
			if r == '/' && !isCurrency {
				state = 2
			} else if r == '.' {
				state = 5
			} else if r == 'k' || r == 'K' || r == 'm' || r == 'M' || r == 'b' || r == 'B' {
				magnitudeSpec = r
				state = 10
			} else if unicode.IsDigit(r) {
				// nop
			} else if unicode.IsPunct(r) {
				inTrailing = true
			} else {
				return NumClassNone
			}
		case 2: // fraction, after slash
			if unicode.IsDigit(r) {
				state = 3
			} else if unicode.IsPunct(r) {
				inTrailing = true
			} else {
				return NumClassNone
			}
		case 3: // fraction, within a run of trailing digits
			if unicode.IsDigit(r) {
				// nop
			} else if unicode.IsPunct(r) {
				inTrailing = true
			} else {
				return NumClassNone
			}
		case 4: // leading digit, after currency marker
			if unicode.IsDigit(r) {
				state = 1
			} else if unicode.IsPunct(r) {
				inTrailing = true
			} else {
				return NumClassNone
			}
		case 5: // after decimal period
			if unicode.IsDigit(r) {
				// nop
			} else if unicode.IsPunct(r) {
				inTrailing = true
			} else {
				return NumClassNone
			}
		case 10: // after magniture specifier ($1M)
			if r == magnitudeSpec {
				state = 11
			} else if unicode.IsPunct(r) {
				inTrailing = true
			} else {
				return NumClassNone
			}
		case 11: // after a double magniture specifier ($1MM)
			if unicode.IsPunct(r) {
				inTrailing = true
			} else {
				return NumClassNone
			}
		}
	}
	if isCurrency {
		if state != 4 {
			return NumClassCurrency
		} else {
			return NumClassNone // just a lone $ character
		}
	}
	if state == 3 {
		return NumClassFraction
	} else if state == 1 {
		return NumClassInteger
	} else if state == 5 {
		return NumClassFloat
	} else {
		return NumClassNone
	}
}
