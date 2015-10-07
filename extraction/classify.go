package extraction

import (
	"fmt"
	"strings"
	"unicode"
)

var ignoredWords = makeTable(`
    from to on in of with
    a an the
    him her them
    our
    out
    all
    this that
    not
    be is are
`)

type unwantedInformationRec struct {
	label          string
	threshold      int
	reqPrimaries   int
	primaryTable   map[string]bool
	secondaryTable map[string]bool
}

var unwantedRecs = []unwantedInformationRec{
	{
		label:        "trailing line",
		threshold:    30,
		reqPrimaries: 1,
		primaryTable: makeTable(`
            contact twitter check podcast download
            rights reserved published copyright associated press
        `),
		secondaryTable: makeTable(`
            free tigers xtra latest material
            may
        `),
	},
	{
		// Anyone with details on what happened can call the Auburn Hills Police Department at 248-370-9444.
		label:        "police",
		threshold:    30,
		reqPrimaries: 3,
		primaryTable: makeTable(`
            happened details call police department
        `),
		secondaryTable: makeTable(`
            anyone can what
            auburn hills
        `),
	},
}

func IsWorthyParagraph(text string) (bool, string) {
	words := SplitWords(text)
	lowerWords := stringsToLower(words)
	total := len(words)

	if total == 0 {
		return false, "empty"
	}

	irregulars := countIrregulars(words)
	ignored := countInTable(ignoredWords, lowerWords)

	rationales := make([]string, 0, len(unwantedRecs))

	for _, unwanted := range unwantedRecs {
		primary := countInTable(unwanted.primaryTable, lowerWords)
		secondary := countInTable(unwanted.secondaryTable, lowerWords)

		percentage := 0
		if primary >= unwanted.reqPrimaries {
			percentage = percentageOf(primary+secondary+irregulars, total-ignored)
		}

		ok := (percentage < unwanted.threshold)
		rationale := fmt.Sprintf("%s (%d%% / %d%%: total %d, irregular %d, primary %d, secondary %d)", unwanted.label, percentage, unwanted.threshold, total, irregulars, primary, secondary)

		if !ok {
			return false, rationale
		} else {
			rationales = append(rationales, "not "+rationale)
		}
	}

	return true, strings.Join(rationales, "; ")
}

func makeTable(spaceSeparated string) map[string]bool {
	items := strings.Fields(spaceSeparated)

	table := make(map[string]bool, len(items))

	for _, item := range items {
		table[item] = true
	}

	return table
}

func stringsToLower(input []string) []string {
	result := make([]string, 0, len(input))
	for _, el := range input {
		el = strings.Map(unicode.ToLower, el)
		result = append(result, el)
	}
	return result
}

func countIrregulars(words []string) int {
	c := 0
	for _, word := range words {
		if !IsRegularWord(word) {
			c++
		}
	}
	return c
}

func countInTable(table map[string]bool, words []string) int {
	c := 0
	for _, word := range words {
		if table[word] {
			c++
		}
	}
	return c
}

func percentageOf(amount, total int) int {
	if total == 0 {
		return 0
	} else {
		return amount * 100 / total
	}
}
