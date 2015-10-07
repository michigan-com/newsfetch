package extraction

import (
	"bytes"

	gq "github.com/PuerkitoBio/goquery"
	m "github.com/michigan-com/newsfetch/model"
)

func ExtractDataFromHTMLString(html string, url string, includeTitle bool) *m.ExtractedBody {
	doc, err := gq.NewDocumentFromReader(bytes.NewBufferString(html))
	if err != nil {
		return nil
	}

	extracted := extractBodyFromDocument(doc, true, includeTitle)
	for _, recipe := range extracted.RecipeData.Recipes {
		recipe.Url = url
	}

	return extracted
}

func ExtractDataFromHTMLAtURL(url string, includeTitle bool) *m.ExtractedBody {
	// Debugger.Printf("Fetching %s ...\n", url)
	doc, err := gq.NewDocument(url)
	if err != nil {
		return nil
	}

	extracted := extractBodyFromDocument(doc, false, includeTitle)
	for _, recipe := range extracted.RecipeData.Recipes {
		recipe.Url = url
	}

	return extracted
}
