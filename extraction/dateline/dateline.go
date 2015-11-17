package dateline

import (
	"strings"
	"unicode"
)

const dashes = "-—"
const skipped = "," + dashes
const maxAlternativeDatelineWords = 5

func RemoveDateline(text string) string {
	found := false
	for {
		brk := strings.IndexFunc(text, unicode.IsSpace)
		if brk < 0 {
			return text
		}

		word := text[0:brk]
		if word != strings.ToUpper(word) {
			if !found {
				return removeDatelineAlternative(text)
			} else {
				return text
			}
		}

		text = text[brk:]
		found = true

		nxt := strings.IndexFunc(text, isNotSkipped)
		if nxt < 0 {
			return ""
		}
		text = text[nxt:]
	}
}

func removeDatelineAlternative(text string) string {
	pos := strings.IndexAny(text, dashes)
	if pos < 0 {
		return text
	}

	dateline := text[0:pos]
	dlwords := strings.Fields(dateline)
	if len(dlwords) > maxAlternativeDatelineWords {
		return text
	}

	text = text[pos:]

	nxt := strings.IndexFunc(text, isNotSkipped)
	if nxt < 0 {
		return ""
	}
	text = text[nxt:]
	return text
}

func isSkipped(r rune) bool {
	return unicode.IsSpace(r) || strings.ContainsRune(skipped, r)
}
func isNotSkipped(r rune) bool {
	return !isSkipped(r)
}

func RmDateline(text string) string {
	emDashes := []string{"—", "--"}
	words := strings.Fields(text)

	wcount := maxAlternativeDatelineWords
	if len(words) < maxAlternativeDatelineWords {
		wcount = len(words)
	}

	truncate := strings.Join(words[:wcount], " ")

	for _, dash := range emDashes {
		dateline := strings.Index(truncate, dash)
		if dateline >= 0 {
			return strings.Trim(text[dateline+len(dash):], " ")
		}
	}

	return text
}
