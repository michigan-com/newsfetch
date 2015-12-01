package fuzzy

import (
	_ "errors"
	"fmt"
	"strings"
	"unicode"

	stemmer "github.com/kljensen/snowball/english"
	"github.com/michigan-com/newsfetch/extraction/split"
)

func (c *Classifier) Process(input string) *Result {
	r := new(Result)
	r.multiVariantTags = c.multiVariantTags

	r.Words = SplitIntoWords(input)

	tagCount := 10 //len(c.categories)
	r.TagsByName = make(map[string][]Range, tagCount)

	r.TagsByPos = make([]map[string][]Range, len(r.Words))
	for pos, _ := range r.Words {
		r.TagsByPos[pos] = make(map[string][]Range, tagCount)
	}

	r.TagDefs = c.TagDefs
	r.TagDefsByName = c.TagDefsByName

	// handle built-in tags
	for pos, word := range r.Words {
		s := word.Trimmed
		runes := []rune(s)
		if len(runes) == 0 {
			continue
		}

		if unicode.IsUpper(runes[0]) && strings.IndexFunc(s, unicode.IsLower) >= 1 {
			r.AddTag("@cap", pos, 1)
		}

		if runes[0] == '@' {
			r.AddTag("@twitter", pos, 1)
		} else if strings.IndexRune(s, '.') >= 1 && strings.IndexRune(s, '/') >= 1 {
			r.AddTag("@url", pos, 1)
		} else if strings.IndexRune(s, '@') >= 1 && strings.IndexRune(s, '.') >= 1 {
			r.AddTag("@email", pos, 1)
		}

		nc := ClassifyNumber(runes)
		switch nc {
		case NumClassInteger:
			r.AddTag("@integer", pos, 1)
		case NumClassFloat:
			r.AddTag("@float", pos, 1)
		case NumClassFraction:
			r.AddTag("@fraction", pos, 1)
		case NumClassCurrency:
			r.AddTag("@currency-number", pos, 1)
		case NumClassNone:
			if !split.IsRegularWord(word.Trimmed) || IsStopWord(word.Stem) {
				r.AddTag("@s", pos, 1)
			}
		}
	}

	// the core: add tags for matched categories
	for _, category := range c.categories {
		// make a list of skip options (each value in skippable is the
		// number of words that can be skipped, ending at this position)
		skippable := make([][]int, len(r.Words))
		for wi, _ := range r.Words {
			for _, tag := range category.skippableTags {
				for _, tagging := range r.TagsByPos[wi][tag] {
					if !intSliceContainsValue(skippable[wi], tagging.Len) {
						last := wi + tagging.Len - 1
						skippable[last] = append(skippable[last], tagging.Len)
					}
				}
			}
		}

		for _, scheme := range category.schemes {
			matchScheme(r, category.tag, scheme.requirements, skippable, category.skipBefore, category.skipAfter)
		}
	}

	return r
}

func intSliceContainsValue(list []int, value int) bool {
	for _, el := range list {
		if el == value {
			return true
		}
	}
	return false
}

// a DP (dynamic programming)-based requirements matcher
//
// parameter 1: number of requirements matched
// (handled implicitly via cur/prev)
//
// parameter 2: index of the last word in a subsequence matching param#1 requirements
// (prev[wi], cur[wi])
//
// Unfortunately, I don't know how to write easy-to-understand algorithmic
// code. This one is fairly standard if you saw your share of DP algorithms.
//
func matchScheme(r *Result, tag string, reqs []Requirement, skippable [][]int, skipBefore, skipAfter bool) {
	matchLogging := false
	// matchLogging = (tag == "@colors")

	words := r.Words
	wc := len(words) + 1
	rc := len(reqs)

	if matchLogging {
		print("----------------------------------------\n")
		for wi := 1; wi < wc; wi++ {
			print("w[", wi, "]=", words[wi-1].Trimmed, " ")
		}
		print("\n")
		for ri, req := range reqs {
			print("req[", ri, "]=", req.String(), " ")
		}
		print("\n")
		print("\n")
	}

	canBeInitial := make([]bool, len(reqs))
	for ri := range reqs {
		if ri == 0 {
			canBeInitial[ri] = true
		} else {
			canBeInitial[ri] = canBeInitial[ri-1] && reqs[ri-1].optional
		}
	}

	prev := make([]int, wc)
	prevDisallowFin := make([]bool, wc)

	// handle skippable words that start a match
	if skipBefore {
		for last := 1; last < wc; last++ {
			for _, skipLen := range skippable[last-1] {
				first := last - skipLen + 1
				if matchLogging {
					print("skip(", first, "..", last, ")\n")
				}

				if prev[first-1] != 0 {
					first = prev[first-1]
				}
				if prev[last] == 0 || prev[last] > first {
					prev[last] = first
				}
			}
		}
	}

	for ri, req := range reqs {
		if matchLogging {
			print("ri=", ri)
			for wi := 0; wi < wc; wi++ {
				print(" prev[", wi, "]=", prev[wi])
			}
			print("\n")
		}

		cur := make([]int, wc)
		curDisallowFin := make([]bool, wc)

		for wi := 1; wi < wc; wi++ {
			first := 0
			if prev[wi-1] != 0 {
				first = prev[wi-1]
			} else if canBeInitial[ri] {
				first = wi
			}
			if first != 0 {
				for _, matchLen := range matchReq(r, req, wi-1) {
					last := wi + matchLen - 1
					if cur[last] == 0 || cur[last] > first {
						cur[last] = first

						if matchLogging {
							print("match(", wi, "..", last, ")")
							for logwi := 0; logwi < wc; logwi++ {
								print(" cur[", logwi, "]=", cur[logwi])
							}
							print("\n")
						}
					}
				}
			}

			if req.repeating {
				if first := cur[wi-1]; first != 0 {
					for _, matchLen := range matchReq(r, req, wi-1) {
						last := wi + matchLen - 1
						if cur[last] == 0 || cur[last] > first {
							cur[last] = first

							if matchLogging {
								print("rep-match(", wi, "..", last, ")")
								for logwi := 0; logwi < wc; logwi++ {
									print(" cur[", logwi, "]=", cur[logwi])
								}
								print("\n")
							}
						}
					}
				}
			}

			if req.optional {
				first := prev[wi]
				if first != 0 {
					if matchLogging {
						print("skip-opt(#", ri, ") first=", first, "\n")
					}
					if cur[wi] == 0 || cur[wi] > first {
						cur[wi] = first
					}
				}
			}

			// handle skippable words that go after ri'th matched requirement
			disallow := (ri == rc-1 && !skipAfter)
			for _, skipLen := range skippable[wi-1] {
				first := wi - skipLen + 1
				if start := cur[first-1]; start != 0 {
					if cur[wi] == 0 || cur[wi] > first {
						if matchLogging {
							print("skip(", first, "..", wi, ") cur[", first-1, "]=", cur[first-1], "\n")
						}
						cur[wi] = start
						curDisallowFin[wi] = curDisallowFin[wi] || disallow
					}
				}
			}
		}

		prev = cur
		prevDisallowFin = curDisallowFin
	}

	for wi := 1; wi < wc; wi++ {
		first := prev[wi]
		if first == 0 {
			continue
		}
		first--
		last := wi - 1

		if prevDisallowFin[wi] {
			if matchLogging {
				print("disallowed(", tag, ": ", first, "..", last, ")\n")
			}
		} else {
			if matchLogging {
				print("tag(", tag, ": ", first, "..", last, ")\n")
			}

			r.AddTag(tag, first, last-first+1)
		}
	}
}

// returns the lengths of the matches
func matchReq(r *Result, req Requirement, wi int) []int {
	switch req.typ {
	case ReqLiteral:
		if req.stem == r.Words[wi].Normalized {
			return []int{1}
		} else {
			return nil
		}

	case ReqStem:
		if req.stem == r.Words[wi].Stem {
			return []int{1}
		} else {
			return nil
		}

	case ReqTag:
		var result []int
		for _, tagging := range r.TagsByPos[wi][req.tag] {
			if !intSliceContainsValue(result, tagging.Len) {
				result = append(result, tagging.Len)
			}
		}
		return result

	default:
		panic("Unknown requirement type")
	}
}

func SplitIntoWords(input string) []Word {
	fields := strings.Fields(input)

	// split words like '2-point', '1-inch' into two words ('2', 'pound')
	mapped := make([]string, 0, len(fields))
	for _, word := range fields {
		trimmed := strings.TrimFunc(word, split.IsTrimmableRune)
		if dash := strings.IndexRune(trimmed, '-'); dash >= 0 {
			left := trimmed[:dash]
			right := trimmed[dash+1:]

			// fmt.Printf("SplitIntoWords dashed %#v, left = %#v, right = %#v, reg(left) = %#v, reg(right) = %#v, num(left) = %#v\n", word, left, right, split.IsRegularWord(left), split.IsRegularWord(right), ClassifyNumberString(left))

			if !split.IsRegularWord(left) && split.IsRegularWord(right) && ClassifyNumberString(left) != NumClassNone {
				idx := strings.Index(word, left)
				if idx < 0 {
					panic(fmt.Sprintf("Cannot find %#v in %#v", left, word))
				}
				dash := idx + len(left)
				fullLeft := word[:dash]
				fullRight := word[dash+1:]

				mapped = append(mapped, fullLeft)
				mapped = append(mapped, fullRight)
				continue
			}
		}
		mapped = append(mapped, word)
	}

	words := make([]Word, 0, len(mapped))
	for _, word := range mapped {
		trimmed := strings.TrimFunc(word, split.IsTrimmableRune)
		trimmed = cleanupWord(trimmed)
		if len(trimmed) == 0 {
			continue
		}

		norm := Normalize(trimmed)

		stem := Stem(trimmed)
		// println(trimmed, "=>", stem)
		words = append(words, Word{Raw: word, Trimmed: trimmed, Stem: stem, Normalized: norm})
	}
	return words
}

func cleanupWord(s string) string {
	s = strings.Replace(s, "½", "1/2", -1)
	s = strings.Replace(s, "¼", "1/4", -1)
	s = strings.Replace(s, "¾", "3/4", -1)
	return s
}

func Normalize(s string) string {
	return strings.ToLower(s)
}

func Stem(s string) string {
	return stemmer.Stem(s, false)
}

// handle cases where SplitIntoWords mutates the words
func CanonicalString(input string) string {
	words := SplitIntoWords(input)
	list := make([]string, 0, len(words))
	for _, word := range words {
		list = append(list, word.Raw)
	}
	return strings.Join(list, " ")
}

// from https://github.com/kljensen/snowball/blob/master/english/common.go
func IsStopWord(word string) bool {
	switch word {
	case "a", "about", "above", "after", "again", "against", "all", "am", "an",
		"and", "any", "are", "as", "at", "be", "because", "been", "before",
		"being", "below", "between", "both", "but", "by", "can", "did", "do",
		"does", "doing", "don", "down", "during", "each", "few", "for", "from",
		"further", "had", "has", "have", "having", "he", "her", "here", "hers",
		"herself", "him", "himself", "his", "how", "i", "if", "in", "into", "is",
		"it", "its", "itself", "just", "me", "more", "most", "my", "myself",
		"no", "nor", "not", "now", "of", "off", "on", "once", "only", "or",
		"other", "our", "ours", "ourselves", "out", "over", "own", "s", "same",
		"she", "should", "so", "some", "such", "t", "than", "that", "the", "their",
		"theirs", "them", "themselves", "then", "there", "these", "they",
		"this", "those", "through", "to", "too", "under", "until", "up",
		"very", "was", "we", "were", "what", "when", "where", "which", "while",
		"who", "whom", "why", "will", "with", "you", "your", "yours", "yourself",
		"yourselves":
		return true
	}
	return false
}
