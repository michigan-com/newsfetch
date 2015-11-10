package classify

import (
	"fmt"
	"strings"

	"github.com/michigan-com/newsfetch/extraction/split"
)

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
	words := split.SplitWords(text)
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
